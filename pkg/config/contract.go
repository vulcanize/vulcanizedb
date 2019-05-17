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
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
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

	// Map of contract address to slice of methods
	// If any methods are listed in the slice only those will be polled
	// Otherwise no methods will be polled
	Methods map[string][]string

	// Map of contract address to slice of event arguments to filter for
	// If arguments are provided then only events which emit those arguments are watched
	// Otherwise arguments are not filtered on events
	EventArgs map[string][]string

	// Map of contract address to slice of method arguments to limit polling to
	// If arguments are provided then only those arguments are allowed as arguments in method polling
	// Otherwise any argument of the right type seen emitted from events at that contract will be used in method polling
	MethodArgs map[string][]string

	// Map of contract address to their starting block
	StartingBlocks map[string]int64

	// Map of contract address to whether or not to pipe method polling results forward into subsequent method calls
	Piping map[string]bool
}

func (contractConfig *ContractConfig) PrepConfig() {
	addrs := viper.GetStringSlice("contract.addresses")
	contractConfig.Network = viper.GetString("contract.network")
	contractConfig.Addresses = make(map[string]bool, len(addrs))
	contractConfig.Abis = make(map[string]string, len(addrs))
	contractConfig.Methods = make(map[string][]string, len(addrs))
	contractConfig.Events = make(map[string][]string, len(addrs))
	contractConfig.MethodArgs = make(map[string][]string, len(addrs))
	contractConfig.EventArgs = make(map[string][]string, len(addrs))
	contractConfig.StartingBlocks = make(map[string]int64, len(addrs))
	contractConfig.Piping = make(map[string]bool, len(addrs))
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
			if _, abiErr := geth.ParseAbi(abi); abiErr != nil {
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

		// Get and check methods
		methods := make([]string, 0)
		methodsInterface, methodsOK := transformer["methods"]
		if !methodsOK {
			log.Warnf("contract %s not configured with a list of methods to poll, will not poll any methods\r\n", addr)
			methods = []string{}
		} else {
			methodsI, methodsOK := methodsInterface.([]interface{})
			if !methodsOK {
				log.Fatal(addr, "transformer `methods` not of type []string\r\n")
			}
			for _, strI := range methodsI {
				str, strOK := strI.(string)
				if !strOK {
					log.Fatal(addr, "transformer `methods` not of type []string\r\n")
				}
				methods = append(methods, str)
			}
		}
		contractConfig.Methods[strings.ToLower(addr)] = methods

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

		// Get and check methodArgs
		methodArgs := make([]string, 0)
		methodArgsInterface, methodArgsOK := transformer["methodArgs"]
		if !methodArgsOK {
			log.Warnf("contract %s not configured with a list of method argument values to poll with, will poll methods with all available arguments\r\n", addr)
			methodArgs = []string{}
		} else {
			methodArgsI, methodArgsOK := methodArgsInterface.([]interface{})
			if !methodArgsOK {
				log.Fatal(addr, "transformer `methodArgs` not of type []string\r\n")
			}
			for _, strI := range methodArgsI {
				str, strOK := strI.(string)
				if !strOK {
					log.Fatal(addr, "transformer `methodArgs` not of type []string\r\n")
				}
				methodArgs = append(methodArgs, str)
			}
		}
		contractConfig.MethodArgs[strings.ToLower(addr)] = methodArgs

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

		// Get pipping
		var piping bool
		_, pipeOK := transformer["piping"]
		if !pipeOK {
			log.Warnf("contract %s does not have its `piping` set, by default piping is turned off\r\n", addr)
			piping = false
		} else {
			pipingInterface := transformer["piping"]
			piping, pipeOK = pipingInterface.(bool)
			if !pipeOK {
				log.Fatal(addr, "transformer `piping` not of type bool\r\n")
			}
		}
		contractConfig.Piping[strings.ToLower(addr)] = piping
	}
}
