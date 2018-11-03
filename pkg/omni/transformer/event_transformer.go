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

// Top-level object similar to generator
// but attempts to solve problem without
// automated code generation

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
	fetcher.Fetcher     // Fetches data from contract methods

	// Store contract info as mapping to contract address
	ContractInfo map[string]types.ContractInfo

	// Subset of events of interest, stored as map of contract address to events
	// By default this
	sets map[string][]string
}

// Transformer takes in config for blockchain, database, and network id
func NewTransformer(c *types.Config) (t *eventTransformer) {
	t.Parser = parser.NewParser(c.Network)
	t.Retriever = retriever.NewRetriever(c.DB)
	t.Fetcher = fetcher.NewFetcher(c.BC)
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
		strName, err1 := t.Fetcher.FetchString("name", t.Parser.GetAbi(), contractAddr, -1, nil)
		if err1 != nil || strName == "" {
			hashName, err2 := t.Fetcher.FetchHash("name", t.Parser.GetAbi(), contractAddr, -1, nil)
			if err2 != nil || hashName.String() == "" {
				return errors.New(fmt.Sprintf("fetching string: %s and hash: %s names failed\r\nerr1: %v\r\nerr2: %v\r\n ", strName, hashName, err1, err2))
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
			Abi:           t.Parser.GetAbi(),
			ParsedAbi:     t.Parser.GetParsedAbi(),
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

func (tr eventTransformer) Execute() error {
	for _, contract := range tr.ContractInfo {

		c := converter.NewConverter(contract)

		for eventName, filter := range contract.Filters {
			watchedEvents, err := tr.GetWatchedEvents(eventName)
			if err != nil {
				log.Println(fmt.Sprintf("Error fetching events for %s:", filter.Name), err)
				return err
			}

			for _, we := range watchedEvents {
				err = c.Convert(*we, contract.Events[eventName])
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
