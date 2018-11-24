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
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/converter"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/parser"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/poller"
	srep "github.com/vulcanize/vulcanizedb/pkg/omni/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

// Requires a light synced vDB (headers) and a running eth node (or infura)
type transformer struct {
	// Database interfaces
	srep.EventRepository        // Holds transformed watched event log data
	repository.HeaderRepository // Interface for interaction with header repositories

	// Pre-processing interfaces
	parser.Parser            // Parses events and methods out of contract abi fetched using contract address
	retriever.BlockRetriever // Retrieves first block for contract and current block height

	// Processing interfaces
	fetcher.Fetcher     // Fetches event logs, using header hashes
	converter.Converter // Converts watched event logs into custom log
	poller.Poller       // Polls methods using contract's token holder addresses and persists them using method datastore

	// Ethereum network name; default "" is mainnet
	Network string

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
func NewTransformer(network string, bc core.BlockChain, db *postgres.DB) *transformer {

	return &transformer{
		Poller:           poller.NewPoller(bc, db, types.LightSync),
		Fetcher:          fetcher.NewFetcher(bc),
		Parser:           parser.NewParser(network),
		HeaderRepository: repository.NewHeaderRepository(db),
		BlockRetriever:   retriever.NewBlockRetriever(db),
		Converter:        converter.NewConverter(&contract.Contract{}),
		Contracts:        map[string]*contract.Contract{},
		EventRepository:  srep.NewEventRepository(db, types.LightSync),
		WatchedEvents:    map[string][]string{},
		WantedMethods:    map[string][]string{},
		ContractRanges:   map[string][2]int64{},
		EventAddrs:       map[string][]string{},
		MethodAddrs:      map[string][]string{},
	}
}

// Use after creating and setting transformer
// Loops over all of the addr => filter sets
// Uses parser to pull event info from abi
// Use this info to generate event filters
func (tr *transformer) Init() error {
	// Iterate through all internal contract addresses
	for contractAddr, subset := range tr.WatchedEvents {
		// Get Abi
		err := tr.Parser.Parse(contractAddr)
		if err != nil {
			return err
		}

		// Get first block and most recent block number in the header repo
		firstBlock, err := tr.BlockRetriever.RetrieveFirstBlock()
		if err != nil {
			return err
		}
		lastBlock, err := tr.BlockRetriever.RetrieveMostRecentBlock()
		if err != nil {
			return err
		}

		// Set to specified range if it falls within the bounds
		if firstBlock < tr.ContractRanges[contractAddr][0] {
			firstBlock = tr.ContractRanges[contractAddr][0]
		}
		if lastBlock > tr.ContractRanges[contractAddr][1] && tr.ContractRanges[contractAddr][1] > firstBlock {
			lastBlock = tr.ContractRanges[contractAddr][1]
		}

		// Get contract name if it has one
		var name = new(string)
		tr.FetchContractData(tr.Abi(), contractAddr, "name", nil, &name, lastBlock)

		// Remove any potential accidental duplicate inputs in filter addresses
		EventAddrs := map[string]bool{}
		for _, addr := range tr.EventAddrs[contractAddr] {
			EventAddrs[addr] = true
		}
		MethodAddrs := map[string]bool{}
		for _, addr := range tr.MethodAddrs[contractAddr] {
			MethodAddrs[addr] = true
		}

		// Aggregate info into contract object
		info := &contract.Contract{
			Name:           *name,
			Network:        tr.Network,
			Address:        contractAddr,
			Abi:            tr.Abi(),
			ParsedAbi:      tr.ParsedAbi(),
			StartingBlock:  firstBlock,
			LastBlock:      lastBlock,
			Events:         tr.GetEvents(subset),
			Methods:        tr.GetAddrMethods(tr.WantedMethods[contractAddr]),
			EventAddrs:     EventAddrs,
			MethodAddrs:    MethodAddrs,
			TknHolderAddrs: map[string]bool{},
		}

		// Store contract info for execution
		tr.Contracts[contractAddr] = info
	}

	return nil
}

func (tr *transformer) Execute() error {
	if len(tr.Contracts) == 0 {
		return errors.New("error: transformer has no initialized contracts")
	}
	// Iterate through all internal contracts
	for _, con := range tr.Contracts {
		// Update converter with current contract
		tr.Update(con)

		// Iterate through events
		for _, event := range con.Events {
			// Filter using the event signature
			topics := [][]common.Hash{{common.HexToHash(helpers.GenerateSignature(event.Sig()))}}

			// Generate eventID and use it to create a checked_header column if one does not already exist
			eventId := strings.ToLower(event.Name + "_" + con.Address)
			if err := tr.AddCheckColumn(eventId); err != nil {
				return err
			}

			// Find unchecked headers for this event
			missingHeaders, err := tr.MissingHeaders(con.StartingBlock, con.LastBlock, eventId)
			if err != nil {
				return err
			}

			// Iterate over headers
			for _, header := range missingHeaders {
				// And fetch event logs using the header, contract address, and topics filter
				logs, err := tr.FetchLogs([]string{con.Address}, topics, header)
				if err != nil {
					return err
				}

				// Mark the header checked for this eventID and continue to next iteration if no logs are found
				if len(logs) < 1 {
					err = tr.MarkHeaderChecked(header.Id, eventId)
					if err != nil {
						return err
					}
					continue
				}

				// Convert logs into custom type
				convertedLogs, err := tr.Convert(logs, event, header.Id)
				if err != nil {
					return err
				}
				if len(convertedLogs) < 1 {
					continue
				}

				// If logs aren't empty, persist them
				err = tr.PersistLogs(convertedLogs, event, con.Address, con.Name)
				if err != nil {
					return err
				}
			}
		}
		// After persisting all watched event logs
		// poller polls select contract methods
		// and persists the results into custom pg tables
		if err := tr.PollContract(*con); err != nil {
			return err
		}
	}

	return nil
}

// Used to set which contract addresses and which of their events to watch
func (tr *transformer) SetEvents(contractAddr string, filterSet []string) {
	tr.WatchedEvents[contractAddr] = filterSet
}

// Used to set subset of account addresses to watch events for
func (tr *transformer) SetEventAddrs(contractAddr string, filterSet []string) {
	tr.EventAddrs[contractAddr] = filterSet
}

// Used to set which contract addresses and which of their methods to call
func (tr *transformer) SetMethods(contractAddr string, filterSet []string) {
	tr.WantedMethods[contractAddr] = filterSet
}

// Used to set subset of account addresses to poll methods on
func (tr *transformer) SetMethodAddrs(contractAddr string, filterSet []string) {
	tr.MethodAddrs[contractAddr] = filterSet
}

// Used to set the block range to watch for a given address
func (tr *transformer) SetRange(contractAddr string, rng [2]int64) {
	tr.ContractRanges[contractAddr] = rng
}
