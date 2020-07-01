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

package config

import (
	"github.com/makerdao/vulcanizedb/pkg/eth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

// Config struct for generic contract transformer
type ContractConfig struct {
	// Name for the transformer
	Name string

	// Ethereum network name; default "" is mainnet
	Network string

	// List of contract addresses (map to ensure no duplicates)
	Addresses map[string]bool

	// Map of contract address to abi
	// If an address has no associated abi the parser will attempt to fetch one from etherscan
	Abis map[string]string

	// Map of contract address to slice of events
	// Used to set which addresses to watch
	// If any events are listed in the slice only those will be watched
	// Otherwise all events in the contract ABI are watched
	Events map[string][]string

	// Map of contract address to slice of event arguments to filter for
	// If arguments are provided then only events which emit those arguments are watched
	// Otherwise arguments are not filtered on events
	EventArgs map[string][]string

	// Map of contract address to their starting block
	StartingBlocks map[string]int64
}

func (contractConfig *ContractConfig) PrepConfig() {
	addrs := viper.GetStringSlice("contract.addresses")
	contractConfig.Network = viper.GetString("contract.network")
	contractConfig.Addresses = make(map[string]bool, len(addrs))
	contractConfig.Abis = make(map[string]string, len(addrs))
	contractConfig.Events = make(map[string][]string, len(addrs))
	contractConfig.EventArgs = make(map[string][]string, len(addrs))
	contractConfig.StartingBlocks = make(map[string]int64, len(addrs))
	// De-dupe addresses
	for _, addr := range addrs {
		contractConfig.Addresses[strings.ToLower(addr)] = true
	}

	// Iterate over addresses to pull out config info for each contract
	for _, addr := range addrs {
		transformer := viper.GetStringMap("contract." + addr)

		// Get and check abi
		var abi string
		abiInterface, abiOK := transformer["abi"]
		if !abiOK {
			log.Warnf("contract %s not configured with an ABI, will attempt to fetch it from Etherscan\r\n", addr)
		} else {
			abi, abiOK = abiInterface.(string)
			if !abiOK {
				log.Fatal(addr, "transformer `abi` not of type []string")
			}
		}
		if abi != "" {
			if _, abiErr := eth.ParseAbi(abi); abiErr != nil {
				log.Fatal(addr, "transformer `abi` not valid JSON")
			}
		}
		contractConfig.Abis[strings.ToLower(addr)] = abi

		// Get and check events
		events := make([]string, 0)
		eventsInterface, eventsOK := transformer["events"]
		if !eventsOK {
			log.Warnf("contract %s not configured with a list of events to watch, will watch all events\r\n", addr)
			events = []string{}
		} else {
			eventsI, eventsOK := eventsInterface.([]interface{})
			if !eventsOK {
				log.Fatal(addr, "transformer `events` not of type []string\r\n")
			}
			for _, strI := range eventsI {
				str, strOK := strI.(string)
				if !strOK {
					log.Fatal(addr, "transformer `events` not of type []string\r\n")
				}
				events = append(events, str)
			}
		}
		contractConfig.Events[strings.ToLower(addr)] = events

		// Get and check eventArgs
		eventArgs := make([]string, 0)
		eventArgsInterface, eventArgsOK := transformer["eventArgs"]
		if !eventArgsOK {
			log.Warnf("contract %s not configured with a list of event arguments to filter for, will not filter events for specific emitted values\r\n", addr)
			eventArgs = []string{}
		} else {
			eventArgsI, eventArgsOK := eventArgsInterface.([]interface{})
			if !eventArgsOK {
				log.Fatal(addr, "transformer `eventArgs` not of type []string\r\n")
			}
			for _, strI := range eventArgsI {
				str, strOK := strI.(string)
				if !strOK {
					log.Fatal(addr, "transformer `eventArgs` not of type []string\r\n")
				}
				eventArgs = append(eventArgs, str)
			}
		}
		contractConfig.EventArgs[strings.ToLower(addr)] = eventArgs

		// Get and check startingBlock
		startInterface, startOK := transformer["startingblock"]
		if !startOK {
			log.Fatal(addr, "transformer config is missing `startingBlock` value\r\n")
		}
		start, startOK := startInterface.(int64)
		if !startOK {
			log.Fatal(addr, "transformer `startingBlock` not of type int\r\n")
		}
		contractConfig.StartingBlocks[strings.ToLower(addr)] = start
	}
}
