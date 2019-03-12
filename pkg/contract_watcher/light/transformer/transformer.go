// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/converter"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/parser"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/poller"
	srep "github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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

	// Store contract configuration information
	Config config.ContractConfig

	// Store contract info as mapping to contract address
	Contracts map[string]*contract.Contract

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
func NewTransformer(con config.ContractConfig, bc core.BlockChain, db *postgres.DB) *transformer {

	return &transformer{
		Poller:           poller.NewPoller(bc, db, types.LightSync),
		Fetcher:          fetcher.NewFetcher(bc),
		Parser:           parser.NewParser(con.Network),
		HeaderRepository: repository.NewHeaderRepository(db),
		BlockRetriever:   retriever.NewBlockRetriever(db),
		Converter:        converter.NewConverter(&contract.Contract{}),
		Contracts:        map[string]*contract.Contract{},
		EventRepository:  srep.NewEventRepository(db, types.LightSync),
		Config:           con,
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
	for contractAddr := range tr.Config.Addresses {
		// Configure Abi
		if tr.Config.Abis[contractAddr] == "" {
			// If no abi is given in the config, this method will try fetching from internal look-up table and etherscan
			err := tr.Parser.Parse(contractAddr)
			if err != nil {
				return err
			}
		} else {
			// If we have an abi from the config, load that into the parser
			err := tr.Parser.ParseAbiStr(tr.Config.Abis[contractAddr])
			if err != nil {
				return err
			}
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
		if firstBlock < tr.Config.StartingBlocks[contractAddr] {
			firstBlock = tr.Config.StartingBlocks[contractAddr]
		}

		// Get contract name if it has one
		var name = new(string)
		tr.FetchContractData(tr.Abi(), contractAddr, "name", nil, &name, lastBlock)

		// Remove any potential accidental duplicate inputs
		eventArgs := map[string]bool{}
		for _, arg := range tr.Config.EventArgs[contractAddr] {
			eventArgs[arg] = true
		}
		methodArgs := map[string]bool{}
		for _, arg := range tr.Config.MethodArgs[contractAddr] {
			methodArgs[arg] = true
		}

		// Aggregate info into contract object and store for execution
		con := contract.Contract{
			Name:          *name,
			Network:       tr.Config.Network,
			Address:       contractAddr,
			Abi:           tr.Parser.Abi(),
			ParsedAbi:     tr.Parser.ParsedAbi(),
			StartingBlock: firstBlock,
			LastBlock:     -1,
			Events:        tr.Parser.GetEvents(tr.Config.Events[contractAddr]),
			Methods:       tr.Parser.GetSelectMethods(tr.Config.Methods[contractAddr]),
			FilterArgs:    eventArgs,
			MethodArgs:    methodArgs,
			Piping:        tr.Config.Piping[contractAddr],
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

		tr.start = header.BlockNumber + 1
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

func (tr *transformer) GetConfig() config.ContractConfig {
	return tr.Config
}
