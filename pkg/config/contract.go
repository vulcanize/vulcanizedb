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

func (oc *ContractConfig) PrepConfig() {
	addrs := viper.GetStringSlice("contract.addresses")
	oc.Network = viper.GetString("contract.network")
	oc.Addresses = make(map[string]bool, len(addrs))
	oc.Abis = make(map[string]string, len(addrs))
	oc.Methods = make(map[string][]string, len(addrs))
	oc.EventArgs = make(map[string][]string, len(addrs))
	oc.MethodArgs = make(map[string][]string, len(addrs))
	oc.EventArgs = make(map[string][]string, len(addrs))
	oc.StartingBlocks = make(map[string]int64, len(addrs))
	oc.Piping = make(map[string]bool, len(addrs))
	// De-dupe addresses
	for _, addr := range addrs {
		oc.Addresses[strings.ToLower(addr)] = true
	}

	// Iterate over addresses to pull out config info for each contract
	for _, addr := range addrs {
		transformer := viper.GetStringMap("contract." + addr)

		// Get and check abi
		abi, abiOK := transformer["abi"]
		if !abiOK || abi == nil {
			log.Fatal(addr, "transformer config is missing `abi` value")
		}
		abiRef, abiOK := abi.(string)
		if !abiOK {
			log.Fatal(addr, "transformer `events` not of type []string")
		}
		oc.Abis[strings.ToLower(addr)] = abiRef

		// Get and check events
		events, eventsOK := transformer["events"]
		if !eventsOK || events == nil {
			log.Fatal(addr, "transformer config is missing `events` value")
		}
		eventsRef, eventsOK := events.([]string)
		if !eventsOK {
			log.Fatal(addr, "transformer `events` not of type []string")
		}
		if eventsRef == nil {
			eventsRef = []string{}
		}
		oc.Events[strings.ToLower(addr)] = eventsRef

		// Get and check methods
		methods, methodsOK := transformer["methods"]
		if !methodsOK || methods == nil {
			log.Fatal(addr, "transformer config is missing `methods` value")
		}
		methodsRef, methodsOK := methods.([]string)
		if !methodsOK {
			log.Fatal(addr, "transformer `methods` not of type []string")
		}
		if methodsRef == nil {
			methodsRef = []string{}
		}
		oc.Methods[strings.ToLower(addr)] = methodsRef

		// Get and check eventArgs
		eventArgs, eventArgsOK := transformer["eventArgs"]
		if !eventArgsOK || eventArgs == nil {
			log.Fatal(addr, "transformer config is missing `eventArgs` value")
		}
		eventArgsRef, eventArgsOK := eventArgs.([]string)
		if !eventArgsOK {
			log.Fatal(addr, "transformer `eventArgs` not of type []string")
		}
		if eventArgsRef == nil {
			eventArgsRef = []string{}
		}
		oc.EventArgs[strings.ToLower(addr)] = eventArgsRef

		// Get and check methodArgs
		methodArgs, methodArgsOK := transformer["methodArgs"]
		if !methodArgsOK || methodArgs == nil {
			log.Fatal(addr, "transformer config is missing `methodArgs` value")
		}
		methodArgsRef, methodArgsOK := methodArgs.([]string)
		if !methodArgsOK {
			log.Fatal(addr, "transformer `methodArgs` not of type []string")
		}
		if methodArgsRef == nil {
			methodArgsRef = []string{}
		}
		oc.MethodArgs[strings.ToLower(addr)] = methodArgsRef

		// Get and check startingBlock
		start, startOK := transformer["startingBlock"]
		if !startOK || start == nil {
			log.Fatal(addr, "transformer config is missing `startingBlock` value")
		}
		startRef, startOK := start.(int64)
		if !startOK {
			log.Fatal(addr, "transformer `startingBlock` not of type int")
		}
		oc.StartingBlocks[strings.ToLower(addr)] = startRef

		// Get pipping
		pipe, pipeOK := transformer["pipping"]
		if !pipeOK || pipe == nil {
			log.Fatal(addr, "transformer config is missing `pipping` value")
		}
		pipeRef, pipeOK := pipe.(bool)
		if !pipeOK {
			log.Fatal(addr, "transformer `piping` not of type bool")
		}
		oc.Piping[strings.ToLower(addr)] = pipeRef
	}
}
