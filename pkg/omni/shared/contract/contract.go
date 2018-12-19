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

package contract

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

// Contract object to hold our contract data
type Contract struct {
	Name           string                       // Name of the contract
	Address        string                       // Address of the contract
	Network        string                       // Network on which the contract is deployed; default empty "" is Ethereum mainnet
	StartingBlock  int64                        // Starting block of the contract
	LastBlock      int64                        // Most recent block on the network
	Abi            string                       // Abi string
	ParsedAbi      abi.ABI                      // Parsed abi
	Events         map[string]types.Event       // List of events to watch
	Methods        []types.Method               // List of methods to poll
	Filters        map[string]filters.LogFilter // Map of event filters to their event names; used only for full sync watcher
	FilterArgs     map[string]bool              // User-input list of values to filter event logs for
	MethodArgs     map[string]bool              // User-input list of values to limit method polling to
	EmittedAddrs   map[interface{}]bool         // List of all unique addresses collected from converted event logs
	EmittedHashes  map[interface{}]bool         // List of all unique hashes collected from converted event logs
	CreateAddrList bool                         // Whether or not to persist address list to postgres
	CreateHashList bool                         // Whether or not to persist hash list to postgres
	Piping         bool                         // Whether or not to pipe method results forward as arguments to subsequent methods
}

// If we will be calling methods that use addr, hash, or byte arrays
// as arguments then we initialize maps to hold these types of values
func (c Contract) Init() *Contract {
	for _, method := range c.Methods {
		for _, arg := range method.Args {
			switch arg.Type.T {
			case abi.AddressTy:
				c.EmittedAddrs = map[interface{}]bool{}
			case abi.HashTy, abi.BytesTy, abi.FixedBytesTy:
				c.EmittedHashes = map[interface{}]bool{}
			default:
			}
		}
	}

	// If we are creating an address list in postgres
	// we initialize the map despite what method call, if any
	if c.CreateAddrList {
		c.EmittedAddrs = map[interface{}]bool{}
	}

	return &c
}

// Use contract info to generate event filters - full sync omni watcher only
func (c *Contract) GenerateFilters() error {
	c.Filters = map[string]filters.LogFilter{}

	for name, event := range c.Events {
		c.Filters[name] = filters.LogFilter{
			Name:      event.Name,
			FromBlock: c.StartingBlock,
			ToBlock:   -1,
			Address:   c.Address,
			Topics:    core.Topics{event.Sig().Hex()},
		}
	}
	// If no filters were generated, throw an error (no point in continuing with this contract)
	if len(c.Filters) == 0 {
		return errors.New("error: no filters created")
	}

	return nil
}

// Returns true if address is in list of arguments to
// filter events for or if no filtering is specified
func (c *Contract) WantedEventArg(arg string) bool {
	if c.FilterArgs == nil {
		return false
	} else if len(c.FilterArgs) == 0 {
		return true
	} else if a, ok := c.FilterArgs[arg]; ok {
		return a
	}

	return false
}

// Returns true if address is in list of arguments to
// poll methods with or if no filtering is specified
func (c *Contract) WantedMethodArg(arg interface{}) bool {
	if c.MethodArgs == nil {
		return false
	} else if len(c.MethodArgs) == 0 {
		return true
	}

	// resolve interface to one of the three types we handle as arguments
	str := StringifyArg(arg)

	// See if it's hex string has been filtered for
	if a, ok := c.MethodArgs[str]; ok {
		return a
	}

	return false
}

// Returns true if any mapping value matches filtered for address or if no filter exists
// Used to check if an event log name-value mapping should be filtered or not
func (c *Contract) PassesEventFilter(args map[string]string) bool {
	for _, arg := range args {
		if c.WantedEventArg(arg) {
			return true
		}
	}

	return false
}

// Add event emitted address to our list if it passes filter and method polling is on
func (c *Contract) AddEmittedAddr(addresses ...interface{}) {
	for _, addr := range addresses {
		if c.WantedMethodArg(addr) && c.Methods != nil {
			c.EmittedAddrs[addr] = true
		}
	}
}

// Add event emitted hash to our list if it passes filter and method polling is on
func (c *Contract) AddEmittedHash(hashes ...interface{}) {
	for _, hash := range hashes {
		if c.WantedMethodArg(hash) && c.Methods != nil {
			c.EmittedHashes[hash] = true
		}
	}
}

func StringifyArg(arg interface{}) (str string) {
	switch arg.(type) {
	case string:
		str = arg.(string)
	case common.Address:
		a := arg.(common.Address)
		str = a.String()
	case common.Hash:
		a := arg.(common.Hash)
		str = a.String()
	case []byte:
		a := arg.([]byte)
		str = hexutil.Encode(a)
	}

	return
}
