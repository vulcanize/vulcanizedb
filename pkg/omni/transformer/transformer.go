// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package transformer

import (
	"errors"
	"fmt"
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
	"github.com/vulcanize/vulcanizedb/pkg/omni/poller"
	"github.com/vulcanize/vulcanizedb/pkg/omni/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/retriever"
)

// Omni transformer
// Used to extract all or a subset of event and method data
// for any contract and persist it to postgres in a manner
// that requires no prior knowledge of the contract
// other than its address and which network it is on
type Transformer interface {
	SetEvents(contractAddr string, filterSet []string)
	SetEventAddrs(contractAddr string, filterSet []string)
	SetMethods(contractAddr string, filterSet []string)
	SetMethodAddrs(contractAddr string, filterSet []string)
	SetRange(contractAddr string, rng [2]int64)
	Init() error
	Execute() error
}

type transformer struct {
	Network string // Ethereum network name; default/"" is mainnet

	// Database interfaces
	datastore.FilterRepository       // Holds watched event log filters
	datastore.WatchedEventRepository // Holds watched event log views, created with the filters
	repository.EventDatastore        // Holds transformed watched event log data

	// Processing interfaces
	parser.Parser              // Parses events out of contract abi fetched with addr
	retriever.BlockRetriever   // Retrieves first block with contract addr referenced
	retriever.AddressRetriever // Retrieves token holder addresses
	converter.Converter        // Converts watched event logs into custom log
	poller.Poller              // Polls methods using contract's token holder addresses

	// Store contract info as mapping to contract address
	Contracts map[string]*contract.Contract

	// Targeted subset of events/methods
	// Stored as map sof contract address to events/method names of interest
	WatchedEvents map[string][]string // Default/empty event list means all are watched
	WantedMethods map[string][]string // Default/empty method list means none are polled

	// Block ranges to watch contracts
	ContractRanges map[string][2]int64

	// Lists of addresses to filter event or method data
	// before persisting; if empty no filter is applied
	EventAddrs  map[string][]string
	MethodAddrs map[string][]string
}

// Transformer takes in config for blockchain, database, and network id
func NewTransformer(network string, BC core.BlockChain, DB *postgres.DB) *transformer {

	return &transformer{
		Poller:                 poller.NewPoller(BC),
		Parser:                 parser.NewParser(network),
		BlockRetriever:         retriever.NewBlockRetriever(DB),
		Converter:              converter.NewConverter(&contract.Contract{}),
		Contracts:              map[string]*contract.Contract{},
		WatchedEventRepository: repositories.WatchedEventRepository{DB: DB},
		FilterRepository:       repositories.FilterRepository{DB: DB},
		EventDatastore:         repository.NewEventDataStore(DB),
		WatchedEvents:          map[string][]string{},
		WantedMethods:          map[string][]string{},
		ContractRanges:         map[string][2]int64{},
		EventAddrs:             map[string][]string{},
		MethodAddrs:            map[string][]string{},
	}
}

// Use after creating and setting transformer
// Loops over all of the addr => filter sets
// Uses parser to pull event info from abi
// Use this info to generate event filters
func (t *transformer) Init() error {

	for contractAddr, subset := range t.WatchedEvents {
		// Get Abi
		err := t.Parser.Parse(contractAddr)
		if err != nil {
			return err
		}

		// Get first block for contract and most recent block for the chain
		firstBlock, err := t.BlockRetriever.RetrieveFirstBlock(contractAddr)
		if err != nil {
			return err
		}
		lastBlock, err := t.BlockRetriever.RetrieveMostRecentBlock()
		if err != nil {
			return err
		}

		// Set to specified range if it falls within the contract's bounds
		if firstBlock < t.ContractRanges[contractAddr][0] {
			firstBlock = t.ContractRanges[contractAddr][0]
		}
		if lastBlock > t.ContractRanges[contractAddr][1] && t.ContractRanges[contractAddr][1] > firstBlock {
			lastBlock = t.ContractRanges[contractAddr][1]
		}

		// Get contract name
		var name = new(string)
		err = t.PollMethod(t.Abi(), contractAddr, "name", nil, &name, lastBlock)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to fetch contract name: %v\r\n", err))
		}

		// Remove any accidental duplicate inputs in filter addresses
		EventAddrs := map[string]bool{}
		for _, addr := range t.EventAddrs[contractAddr] {
			EventAddrs[addr] = true
		}
		MethodAddrs := map[string]bool{}
		for _, addr := range t.MethodAddrs[contractAddr] {
			MethodAddrs[addr] = true
		}

		// Aggregate info into contract object
		info := &contract.Contract{
			Name:           *name,
			Address:        contractAddr,
			Abi:            t.Abi(),
			ParsedAbi:      t.ParsedAbi(),
			StartingBlock:  firstBlock,
			LastBlock:      lastBlock,
			Events:         t.GetEvents(subset),
			Methods:        t.GetAddrMethods(t.WantedMethods[contractAddr]),
			EventAddrs:     EventAddrs,
			MethodAddrs:    MethodAddrs,
			TknHolderAddrs: map[string]bool{},
		}

		// Use info to create filters
		err = info.GenerateFilters()
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
		tr.Converter.Update(con)

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
		err := tr.PersistContractEvents(con)
		if err != nil {
			return err
		}
	}

	return nil
}

// Used to set which contract addresses and which of their events to watch
func (t *transformer) SetEvents(contractAddr string, filterSet []string) {
	t.WatchedEvents[contractAddr] = filterSet
}

// Used to set subset of account addresses to watch events for
func (t *transformer) SetEventAddrs(contractAddr string, filterSet []string) {
	t.EventAddrs[contractAddr] = filterSet
}

// Used to set which contract addresses and which of their methods to call
func (t *transformer) SetMethods(contractAddr string, filterSet []string) {
	t.WantedMethods[contractAddr] = filterSet
}

// Used to set subset of account addresses to poll methods on
func (t *transformer) SetMethodAddrs(contractAddr string, filterSet []string) {
	t.MethodAddrs[contractAddr] = filterSet
}

// Used to set the block range to watch for a given address
func (t *transformer) SetRange(contractAddr string, rng [2]int64) {
	t.ContractRanges[contractAddr] = rng
}
