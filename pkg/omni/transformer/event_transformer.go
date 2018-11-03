// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformer

import (
	"errors"
	"fmt"
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
	"github.com/vulcanize/vulcanizedb/pkg/omni/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Omni event transformer
// Used to extract all or a subset of event data for
// any contract and persist it to postgres in a manner
// that requires no prior knowledge of the contract
// other than its address and which network it is on
type EventTransformer interface {
	Init(contractAddr string) error
}

type eventTransformer struct {
	// Network, database, and blockchain config
	*types.Config

	// Underlying databases
	datastore.WatchedEventRepository
	datastore.FilterRepository
	repository.DataStore

	// Underlying interfaces
	parser.Parser       // Parses events out of contract abi fetched with addr
	retriever.Retriever // Retrieves first block with contract addr referenced
	fetcher.Fetcher     // Fetches data from public contract methods
	converter.Converter // Converts watched event logs into custom log

	// Store contract info as mapping to contract address
	ContractInfo map[string]types.ContractInfo

	// Subset of events of interest, stored as map of contract address to events
	// Default/empty list means all events are considered for that address
	sets map[string][]string
}

// Transformer takes in config for blockchain, database, and network id
func NewTransformer(c *types.Config) (t *eventTransformer) {
	t.Parser = parser.NewParser(c.Network)
	t.Retriever = retriever.NewRetriever(c.DB)
	t.Fetcher = fetcher.NewFetcher(c.BC)
	t.Converter = converter.NewConverter(types.ContractInfo{})
	t.ContractInfo = map[string]types.ContractInfo{}
	t.WatchedEventRepository = repositories.WatchedEventRepository{DB: c.DB}
	t.FilterRepository = repositories.FilterRepository{DB: c.DB}
	t.DataStore = repository.NewDataStore(c.DB)
	t.sets = map[string][]string{}

	return t
}

// Used to set which contract addresses and which of their events to watch
func (t *eventTransformer) Set(contractAddr string, filterSet []string) {
	t.sets[contractAddr] = filterSet
}

// Use after creating and setting transformer
// Loops over all of the addr => filter sets
// Uses parser to pull event info from abi
// Use this info to generate event filters
func (t *eventTransformer) Init() error {

	for contractAddr, subset := range t.sets {
		err := t.Parser.Parse(contractAddr)
		if err != nil {
			return err
		}

		var ctrName string
		strName, err := t.Fetcher.FetchString("name", t.Parser.Abi(), contractAddr, -1, nil)
		if err != nil || strName == "" {
			hashName, err := t.Fetcher.FetchHash("name", t.Parser.Abi(), contractAddr, -1, nil)
			if err != nil || hashName.String() == "" {
				return errors.New("unable to fetch contract name") // provide CLI prompt here for user to input a contract name?
			}
			ctrName = hashName.String()
		} else {
			ctrName = strName
		}

		firstBlock, err := t.Retriever.RetrieveFirstBlock(contractAddr)
		if err != nil {
			return err
		}

		info := types.ContractInfo{
			Name:          ctrName,
			Address:       contractAddr,
			Abi:           t.Parser.Abi(),
			ParsedAbi:     t.Parser.ParsedAbi(),
			StartingBlock: firstBlock,
			Events:        t.Parser.GetEvents(),
			Methods:       t.Parser.GetMethods(),
		}

		info.GenerateFilters(subset)

		for _, filter := range info.Filters {
			t.CreateFilter(filter)
		}

		t.ContractInfo[contractAddr] = info
	}

	return nil
}

// Iterate through contracts, creating a new
// converter for each one and using it to
// convert watched event logs and persist
// them into the postgres db
func (tr eventTransformer) Execute() error {
	for _, contract := range tr.ContractInfo {

		tr.Converter.Update(contract)

		for eventName, filter := range contract.Filters {
			watchedEvents, err := tr.GetWatchedEvents(eventName)
			if err != nil {
				log.Println(fmt.Sprintf("Error fetching events for %s:", filter.Name), err)
				return err
			}

			for _, we := range watchedEvents {
				err = tr.Converter.Convert(*we, contract.Events[eventName])
				if err != nil {
					return err
				}
			}

		}

		err := tr.PersistEvents(contract)
		if err != nil {
			return err
		}
	}

	return nil
}
