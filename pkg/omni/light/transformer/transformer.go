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
	poller.Poller       // Polls methods using arguments collected from events and persists them using a method datastore

	// Ethereum network name; default "" is mainnet
	Network string

	// Store contract info as mapping to contract address
	Contracts map[string]*contract.Contract

	// Targeted subset of events/methods
	// Stored as maps of contract address to events/method names of interest
	WatchedEvents map[string][]string // Default/empty event list means all are watched
	WantedMethods map[string][]string // Default/empty method list means none are polled

	// Starting block number for each contract
	ContractStart map[string]int64

	// Lists of argument values to filter event or
	// method data with; if empty no filter is applied
	EventArgs  map[string][]string
	MethodArgs map[string][]string

	// Whether or not to create a list of emitted address or hashes for the contract in postgres
	CreateAddrList map[string]bool
	CreateHashList map[string]bool

	// Method piping on/off for a contract
	Piping map[string]bool

	// Internally configured transformer variables
	contractAddresses []string            // Holds all contract addresses, for batch fetching of logs
	sortedEventIds    map[string][]string // Map to sort event column ids by contract, for post fetch processing and persisting of logs
	sortedMethodIds   map[string][]string // Map to sort method column ids by contract, for post fetch method polling
	eventIds          []string            // Holds event column ids across all contract, for batch fetching of headers
	eventFilters      []common.Hash       // Holds topic0 hashes across all contracts, for batch fetching of logs
	start             int64               // Hold the lowest starting block and the highest ending block
}

// Order-of-operations:
// 1. Create new transformer
// 2. Load contract addresses and their parameters
// 3. Init
// 4. Execute

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
		ContractStart:    map[string]int64{},
		EventArgs:        map[string][]string{},
		MethodArgs:       map[string][]string{},
		CreateAddrList:   map[string]bool{},
		CreateHashList:   map[string]bool{},
		Piping:           map[string]bool{},
		Network:          network,
	}
}

// Use after creating and setting transformer
// Loops over all of the addr => filter sets
// Uses parser to pull event info from abi
// Use this info to generate event filters
func (tr *transformer) Init() error {
	// Initialize internally configured transformer settings
	tr.contractAddresses = make([]string, 0)       // Holds all contract addresses, for batch fetching of logs
	tr.sortedEventIds = make(map[string][]string)  // Map to sort event column ids by contract, for post fetch processing and persisting of logs
	tr.sortedMethodIds = make(map[string][]string) // Map to sort method column ids by contract, for post fetch method polling
	tr.eventIds = make([]string, 0)                // Holds event column ids across all contract, for batch fetching of headers
	tr.eventFilters = make([]common.Hash, 0)       // Holds topic0 hashes across all contracts, for batch fetching of logs
	tr.start = 100000000000                        // Hold the lowest starting block and the highest ending block

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
		if firstBlock < tr.ContractStart[contractAddr] {
			firstBlock = tr.ContractStart[contractAddr]
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
		con := contract.Contract{
			Name:           *name,
			Network:        tr.Network,
			Address:        contractAddr,
			Abi:            tr.Parser.Abi(),
			ParsedAbi:      tr.Parser.ParsedAbi(),
			StartingBlock:  firstBlock,
			LastBlock:      -1,
			Events:         tr.Parser.GetEvents(subset),
			Methods:        tr.Parser.GetSelectMethods(tr.WantedMethods[contractAddr]),
			FilterArgs:     eventArgs,
			MethodArgs:     methodArgs,
			CreateAddrList: tr.CreateAddrList[contractAddr],
			CreateHashList: tr.CreateHashList[contractAddr],
			Piping:         tr.Piping[contractAddr],
		}.Init()
		tr.Contracts[contractAddr] = con
		tr.contractAddresses = append(tr.contractAddresses, con.Address)

		// Create checked_headers columns for each event id and append to list of all event ids
		tr.sortedEventIds[con.Address] = make([]string, 0, len(con.Events))
		for _, event := range con.Events {
			eventId := strings.ToLower(event.Name + "_" + con.Address)
			err := tr.HeaderRepository.AddCheckColumn(eventId)
			if err != nil {
				return err
			}
			// Keep track of this event id; sorted and unsorted
			tr.sortedEventIds[con.Address] = append(tr.sortedEventIds[con.Address], eventId)
			tr.eventIds = append(tr.eventIds, eventId)
			// Append this event sig to the filters
			tr.eventFilters = append(tr.eventFilters, event.Sig())
		}

		// Create checked_headers columns for each method id and append list of all method ids
		tr.sortedMethodIds[con.Address] = make([]string, 0, len(con.Methods))
		for _, m := range con.Methods {
			methodId := strings.ToLower(m.Name + "_" + con.Address)
			err := tr.HeaderRepository.AddCheckColumn(methodId)
			if err != nil {
				return err
			}
			tr.sortedMethodIds[con.Address] = append(tr.sortedMethodIds[con.Address], methodId)
		}

		// Update start to the lowest block
		if con.StartingBlock < tr.start {
			tr.start = con.StartingBlock
		}
	}

	return nil
}

func (tr *transformer) Execute() error {
	if len(tr.Contracts) == 0 {
		return errors.New("error: transformer has no initialized contracts")
	}

	// Map to sort batch fetched logs by which contract they belong to, for post fetch processing
	sortedLogs := make(map[string][]gethTypes.Log)
	for _, con := range tr.Contracts {
		sortedLogs[con.Address] = []gethTypes.Log{}
	}

	// Find unchecked headers for all events across all contracts; these are returned in asc order
	missingHeaders, err := tr.HeaderRepository.MissingHeadersForAll(tr.start, -1, tr.eventIds)
	if err != nil {
		return err
	}

	// Iterate over headers
	for _, header := range missingHeaders {
		// And fetch all event logs across contracts at this header
		allLogs, err := tr.Fetcher.FetchLogs(tr.contractAddresses, tr.eventFilters, header)
		if err != nil {
			return err
		}

		// If no logs are found mark the header checked for all of these eventIDs
		// and continue to method polling and onto the next iteration
		if len(allLogs) < 1 {
			err = tr.HeaderRepository.MarkHeaderCheckedForAll(header.Id, tr.eventIds)
			if err != nil {
				return err
			}
			err = tr.methodPolling(header, tr.sortedMethodIds)
			if err != nil {
				return err
			}
			continue
		}

		// Sort logs by the contract they belong to
		for _, log := range allLogs {
			addr := strings.ToLower(log.Address.Hex())
			sortedLogs[addr] = append(sortedLogs[addr], log)
		}

		// Process logs for each contract
		for conAddr, logs := range sortedLogs {
			if logs == nil {
				continue
			}
			// Configure converter with this contract
			con := tr.Contracts[conAddr]
			tr.Converter.Update(con)

			// Convert logs into batches of log mappings (eventName => []types.Logs
			convertedLogs, err := tr.Converter.ConvertBatch(logs, con.Events, header.Id)
			if err != nil {
				return err
			}
			// Cycle through each type of event log and persist them
			for eventName, logs := range convertedLogs {
				// If logs for this event are empty, mark them checked at this header and continue
				if len(logs) < 1 {
					eventId := strings.ToLower(eventName + "_" + con.Address)
					err = tr.HeaderRepository.MarkHeaderChecked(header.Id, eventId)
					if err != nil {
						return err
					}
					continue
				}
				// If logs aren't empty, persist them
				// Header is marked checked in the transactions
				err = tr.EventRepository.PersistLogs(logs, con.Events[eventName], con.Address, con.Name)
				if err != nil {
					return err
				}
			}
		}

		// Poll contracts at this block height
		err = tr.methodPolling(header, tr.sortedMethodIds)
		if err != nil {
			return err
		}
	}

	return nil
}

// Used to poll contract methods at a given header
func (tr *transformer) methodPolling(header core.Header, sortedMethodIds map[string][]string) error {
	for _, con := range tr.Contracts {
		// Skip method polling processes if no methods are specified
		// Also don't try to poll methods below this contract's specified starting block
		if len(con.Methods) == 0 || header.BlockNumber < con.StartingBlock {
			continue
		}

		// Poll all methods for this contract at this header
		err := tr.Poller.PollContractAt(*con, header.BlockNumber)
		if err != nil {
			return err
		}

		// Mark this header checked for the methods
		err = tr.HeaderRepository.MarkHeaderCheckedForAll(header.Id, sortedMethodIds[con.Address])
		if err != nil {
			return err
		}
	}

	return nil
}

// Used to set which contract addresses and which of their events to watch
func (tr *transformer) SetEvents(contractAddr string, filterSet []string) {
	tr.WatchedEvents[strings.ToLower(contractAddr)] = filterSet
}

// Used to set subset of account addresses to watch events for
func (tr *transformer) SetEventArgs(contractAddr string, filterSet []string) {
	tr.EventArgs[strings.ToLower(contractAddr)] = filterSet
}

// Used to set which contract addresses and which of their methods to call
func (tr *transformer) SetMethods(contractAddr string, filterSet []string) {
	tr.WantedMethods[strings.ToLower(contractAddr)] = filterSet
}

// Used to set subset of account addresses to poll methods on
func (tr *transformer) SetMethodArgs(contractAddr string, filterSet []string) {
	tr.MethodArgs[strings.ToLower(contractAddr)] = filterSet
}

// Used to set the block range to watch for a given address
func (tr *transformer) SetStartingBlock(contractAddr string, start int64) {
	tr.ContractStart[strings.ToLower(contractAddr)] = start
}

// Used to set whether or not to persist an account address list
func (tr *transformer) SetCreateAddrList(contractAddr string, on bool) {
	tr.CreateAddrList[strings.ToLower(contractAddr)] = on
}

// Used to set whether or not to persist an hash list
func (tr *transformer) SetCreateHashList(contractAddr string, on bool) {
	tr.CreateHashList[strings.ToLower(contractAddr)] = on
}

// Used to turn method piping on for a contract
func (tr *transformer) SetPiping(contractAddr string, on bool) {
	tr.Piping[strings.ToLower(contractAddr)] = on
}
