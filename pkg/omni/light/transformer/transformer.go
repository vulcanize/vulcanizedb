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
	"strings"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"

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
	EventArgs  map[string][]string
	MethodArgs map[string][]string

	// Whether or not to create a list of emitted address or hashes for the contract in postgres
	CreateAddrList map[string]bool
	CreateHashList map[string]bool

	// Method piping on/off for a contract
	Piping map[string]bool
}

// Order-of-operations:
// 1. Create new transformer
// 2. Load contract addresses and their parameters
// 3. Init
// 3. Execute

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
		EventArgs:        map[string][]string{},
		MethodArgs:       map[string][]string{},
		CreateAddrList:   map[string]bool{},
		CreateHashList:   map[string]bool{},
		Piping:           map[string]bool{},
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

		// Remove any potential accidental duplicate inputs in arg filter values
		eventArgs := map[string]bool{}
		for _, arg := range tr.EventArgs[contractAddr] {
			eventArgs[arg] = true
		}
		methodArgs := map[string]bool{}
		for _, arg := range tr.MethodArgs[contractAddr] {
			methodArgs[arg] = true
		}

		// Aggregate info into contract object and store for execution
		tr.Contracts[contractAddr] = contract.Contract{
			Name:           *name,
			Network:        tr.Network,
			Address:        contractAddr,
			Abi:            tr.Abi(),
			ParsedAbi:      tr.ParsedAbi(),
			StartingBlock:  firstBlock,
			LastBlock:      lastBlock,
			Events:         tr.GetEvents(subset),
			Methods:        tr.GetSelectMethods(tr.WantedMethods[contractAddr]),
			FilterArgs:     eventArgs,
			MethodArgs:     methodArgs,
			CreateAddrList: tr.CreateAddrList[contractAddr],
			CreateHashList: tr.CreateHashList[contractAddr],
			Piping:         tr.Piping[contractAddr],
		}.Init()
	}

	return nil
}

func (tr *transformer) Execute() error {
	if len(tr.Contracts) == 0 {
		return errors.New("error: transformer has no initialized contracts")
	}

	cLen := len(tr.Contracts)
	contractAddresses := make([]string, 0, cLen)   // Holds all contract addresses, for batch fetching of logs
	sortedIds := make(map[string][]string)         // Map to sort event column ids by contract, for post fetch processing
	eventIds := make([]string, 0)                  // Holds event column ids across all contract, for batch fetching of headers
	eventFilters := make([]common.Hash, 0)         // Holds topic hashes across all contracts, for batch fetching of logs
	sortedLogs := make(map[string][]gethTypes.Log) // Map to sort batch fetched logs by which contract they belong to, for post fetch processing
	var start, end int64                           // Hold the lowest starting block and the highest ending block
	start = 100000000
	end = -1

	// Cycle through all contracts and extract info needed for fetching and post-processing
	for _, con := range tr.Contracts {
		eLen := len(con.Events)
		sortedLogs[con.Address] = []gethTypes.Log{}
		sortedIds[con.Address] = make([]string, 0, eLen)
		for _, event := range con.Events {
			// Generate eventID and use it to create a checked_header column if one does not already exist
			eventId := strings.ToLower(event.Name + "_" + con.Address)
			err := tr.HeaderRepository.AddCheckColumn(eventId)
			if err != nil {
				return err
			}
			// Keep track of this event id; sorted and unsorted
			sortedIds[con.Address] = append(sortedIds[con.Address], eventId)
			eventIds = append(eventIds, eventId)
			// Append this event sig to the filters
			eventFilters = append(eventFilters, event.Sig())
		}
		contractAddresses = append(contractAddresses, con.Address)
		// Update start to the lowest block and end to the highest block
		if con.StartingBlock < start {
			start = con.StartingBlock
		}
		if con.LastBlock > end {
			end = con.LastBlock
		}
	}

	// Find unchecked headers for all events across all contracts
	missingHeaders, err := tr.HeaderRepository.MissingHeadersForAll(start, end, eventIds)
	if err != nil {
		return err
	}

	// Iterate over headers
	for _, header := range missingHeaders {
		// And fetch all event logs across contracts at this header
		allLogs, err := tr.Fetcher.FetchLogs(contractAddresses, eventFilters, header)
		if err != nil {
			return err
		}

		// Mark the header checked for all of these eventIDs and continue to next iteration if no logs are found
		if len(allLogs) < 1 {
			err = tr.HeaderRepository.MarkHeaderCheckedForAll(header.Id, eventIds)
			if err != nil {
				return err
			}
			continue
		}

		// Sort logs by the contract they belong to
		for _, log := range allLogs {
			sortedLogs[log.Address.Hex()] = append(sortedLogs[log.Address.Hex()], log)
		}

		// Process logs for each contract
		for conAddr, logs := range sortedLogs {
			// Configure converter with this contract
			con := tr.Contracts[conAddr]
			tr.Converter.Update(con)

			// Convert logs into batches of log mappings (event => []types.Log)
			convertedLogs, err := tr.Converter.ConvertBatch(logs, con.Events, header.Id)
			if err != nil {
				return err
			}

			// Cycle through each type of event log and persist them
			for eventName, logs := range convertedLogs {
				// If logs are empty, mark checked
				if len(logs) < 1 {
					eventId := strings.ToLower(eventName + "_" + con.Address)
					err = tr.HeaderRepository.MarkHeaderChecked(header.Id, eventId)
					if err != nil {
						return err
					}
					continue
				}
				// If logs aren't empty, persist them
				// Headers are marked checked in the persistlogs transactions
				err = tr.EventRepository.PersistLogs(logs, con.Events[eventName], con.Address, con.Name)
				if err != nil {
					return err
				}
			}

			// Skip method polling processes if no methods are specified
			if len(con.Methods) == 0 {
				continue
			}

			// Create checked_headers columns for each method id and generate list of all method ids
			methodIds := make([]string, 0, len(con.Methods))
			for _, m := range con.Methods {
				methodId := strings.ToLower(m.Name + "_" + con.Address)
				err = tr.HeaderRepository.AddCheckColumn(methodId)
				if err != nil {
					return err
				}
				methodIds = append(methodIds, methodId)
			}

			// Retrieve headers that have been checked for all of this contract's events but haven not been checked for this contract's methods
			missingHeaders, err = tr.HeaderRepository.MissingMethodsCheckedEventsIntersection(con.StartingBlock, con.LastBlock, methodIds, sortedIds[conAddr])
			if err != nil {
				return err
			}

			// Poll over the missing headers
			for _, header := range missingHeaders {
				err = tr.Poller.PollContractAt(*con, header.BlockNumber)
				if err != nil {
					return err
				}
			}

			// Mark those headers checked for the methods
			err = tr.HeaderRepository.MarkHeadersCheckedForAll(missingHeaders, methodIds)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Used to set which contract addresses and which of their events to watch
func (tr *transformer) SetEvents(contractAddr string, filterSet []string) {
	tr.WatchedEvents[contractAddr] = filterSet
}

// Used to set subset of account addresses to watch events for
func (tr *transformer) SetEventArgs(contractAddr string, filterSet []string) {
	tr.EventArgs[contractAddr] = filterSet
}

// Used to set which contract addresses and which of their methods to call
func (tr *transformer) SetMethods(contractAddr string, filterSet []string) {
	tr.WantedMethods[contractAddr] = filterSet
}

// Used to set subset of account addresses to poll methods on
func (tr *transformer) SetMethodArgs(contractAddr string, filterSet []string) {
	tr.MethodArgs[contractAddr] = filterSet
}

// Used to set the block range to watch for a given address
func (tr *transformer) SetRange(contractAddr string, rng [2]int64) {
	tr.ContractRanges[contractAddr] = rng
}

// Used to set whether or not to persist an account address list
func (tr *transformer) SetCreateAddrList(contractAddr string, on bool) {
	tr.CreateAddrList[contractAddr] = on
}

// Used to set whether or not to persist an hash list
func (tr *transformer) SetCreateHashList(contractAddr string, on bool) {
	tr.CreateHashList[contractAddr] = on
}

// Used to turn method piping on for a contract
func (tr *transformer) SetPiping(contractAddr string, on bool) {
	tr.Piping[contractAddr] = on
}
