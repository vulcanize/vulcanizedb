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

package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/types"
)

// Contract object to hold our contract data
type Contract struct {
	Address       string                 // Address of the contract
	Network       string                 // Network on which the contract is deployed; default empty "" is Ethereum mainnet
	StartingBlock int64                  // Starting block of the contract
	Abi           string                 // Abi string
	ParsedAbi     abi.ABI                // Parsed abi
	Events        map[string]types.Event // List of events to watch
	FilterArgs    map[string]bool        // User-input list of values to filter event logs for
}

// Init initializes a contract object
func (c Contract) Init() *Contract {
	return &c
}

// WantedEventArg returns true if address is in list of arguments to
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

// PassesEventFilter returns true if any mapping value matches filtered for address or if no filter exists
// Used to check if an event log name-value mapping should be filtered or not
func (c *Contract) PassesEventFilter(args map[string]string) bool {
	for _, arg := range args {
		if c.WantedEventArg(arg) {
			return true
		}
	}

	return false
}
