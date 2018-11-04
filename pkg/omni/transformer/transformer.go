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
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
	"github.com/vulcanize/vulcanizedb/pkg/omni/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Omni transformer
// Used to extract all or a subset of event and method data
// for any contract and persist it to postgres in a manner
// that requires no prior knowledge of the contract
// other than its address and which network it is on
type Transformer interface {
	Init(contractAddr string) error
}

type transformer struct {
	// Network, database, and blockchain config
	*types.Config

	// Underlying databases
	datastore.WatchedEventRepository
	datastore.FilterRepository
	repository.DataStore

	// Underlying interfaces
	parser.Parser              // Parses events out of contract abi fetched with addr
	retriever.BlockRetriever   // Retrieves first block with contract addr referenced
	retriever.AddressRetriever // Retrieves token holder addresses
	fetcher.Fetcher            // Fetches data from public contract methods
	converter.Converter        // Converts watched event logs into custom log

	// Store contract info as mapping to contract address
	Contracts map[string]*contract.Contract

	// Targeted subset of events/methods
	// Stored as map of contract address to events/method names of interest
	// Default/empty list means all events/methods are considered for that address
	targetEvents  map[string][]string
	targetMethods map[string][]string
}

// Transformer takes in config for blockchain, database, and network id
func NewTransformer(c *types.Config) *transformer {

	return &transformer{
		Parser:                 parser.NewParser(c.Network),
		BlockRetriever:         retriever.NewBlockRetriever(c.DB),
		Fetcher:                fetcher.NewFetcher(c.BC),
		Converter:              converter.NewConverter(contract.Contract{}),
		Contracts:              map[string]*contract.Contract{},
		WatchedEventRepository: repositories.WatchedEventRepository{DB: c.DB},
		FilterRepository:       repositories.FilterRepository{DB: c.DB},
		DataStore:              repository.NewDataStore(c.DB),
		targetEvents:           map[string][]string{},
		targetMethods:          map[string][]string{},
	}
}

// Used to set which contract addresses and which of their events to watch
func (t *transformer) SetEvents(contractAddr string, filterSet []string) {
	t.targetEvents[contractAddr] = filterSet
}

// Used to set which contract addresses and which of their methods to call
func (t *transformer) SetMethods(contractAddr string, filterSet []string) {
	t.targetMethods[contractAddr] = filterSet
}

// Use after creating and setting transformer
// Loops over all of the addr => filter sets
// Uses parser to pull event info from abi
// Use this info to generate event filters
func (t *transformer) Init() error {

	for contractAddr, subset := range t.targetEvents {
		// Get Abi
		err := t.Parser.Parse(contractAddr)
		if err != nil {
			return err
		}

		// Get first block for contract
		firstBlock, err := t.BlockRetriever.RetrieveFirstBlock(contractAddr)
		if err != nil {
			return err
		}

		// Get most recent block
		lastBlock, err := t.BlockRetriever.RetrieveMostRecentBlock()
		if err != nil {
			return err
		}

		// Get contract name
		var ctrName string // should change this to check for "name" method and its return type in the abi methods before trying to fetch
		strName, err := t.Fetcher.FetchString("name", t.Parser.Abi(), contractAddr, lastBlock, nil)
		if err != nil || strName == "" {
			hashName, err := t.Fetcher.FetchHash("name", t.Parser.Abi(), contractAddr, lastBlock, nil)
			if err != nil || hashName.String() == "" {
				return errors.New(fmt.Sprintf("unable to fetch contract name: %v\r\n", err)) // provide CLI prompt here for user to input a contract name?
			}
			ctrName = hashName.String()
		} else {
			ctrName = strName
		}

		// Aggregate info into contract object
		info := &contract.Contract{
			Name:          ctrName,
			Address:       contractAddr,
			Abi:           t.Parser.Abi(),
			ParsedAbi:     t.Parser.ParsedAbi(),
			StartingBlock: firstBlock,
			Events:        t.Parser.GetEvents(),
			Methods:       t.Parser.GetMethods(),
			Addresses:     map[string]bool{},
		}

		// Use info to create filters
		err = info.GenerateFilters(subset)
		if err != nil {
			return err
		}

		// Iterate over filters and push them to the repo
		for _, filter := range info.Filters {
			t.CreateFilter(filter)
		}

		t.Contracts[contractAddr] = info
	}

	return nil
}

// Iterate through contracts, updating the
// converter with each one and using it to
// convert watched event logs.
// Then persist them into the postgres db
func (tr transformer) Execute() error {
	// Iterate through all internal contracts
	for _, con := range tr.Contracts {

		// Update converter with current contract
		tr.Converter.Update(*con)

		// Iterate through contract filters and get watched event logs
		for eventName, filter := range con.Filters {
			watchedEvents, err := tr.GetWatchedEvents(eventName)
			if err != nil {
				log.Println(fmt.Sprintf("Error fetching events for %s:", filter.Name), err)
				return err
			}

			// Iterate over watched event logs and convert them
			for _, we := range watchedEvents {
				err = tr.Converter.Convert(*we, con.Events[eventName])
				if err != nil {
					return err
				}
			}

		}

		// After converting all logs for events of interest, persist all of the data
		err := tr.PersistEvents(con)
		if err != nil {
			return err
		}
	}

	return nil
}
