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

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/omni/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

type Contract struct {
	Name           string
	Address        string
	StartingBlock  int64
	LastBlock      int64
	Abi            string
	ParsedAbi      abi.ABI
	Events         map[string]*types.Event      // Map of events to their names
	Methods        map[string]*types.Method     // Map of methods to their names
	Filters        map[string]filters.LogFilter // Map of event filters to their names
	EventAddrs     map[string]bool              // User-input list of account addresses to watch events for
	MethodAddrs    map[string]bool              // User-input list of account addresses to poll methods for
	TknHolderAddrs map[string]bool              // List of all contract-associated addresses, populated as events are transformed
}

// Use contract info to generate event filters
func (c *Contract) GenerateFilters() error {
	c.Filters = map[string]filters.LogFilter{}

	for name, event := range c.Events {
		c.Filters[name] = filters.LogFilter{
			Name:      name,
			FromBlock: c.StartingBlock,
			ToBlock:   -1,
			Address:   c.Address,
			Topics:    core.Topics{helpers.GenerateSignature(event.Sig())}, // move generate signatrue to pkg
		}
	}
	// If no filters we generated, throw an error (no point in continuing)
	if len(c.Filters) == 0 {
		return errors.New("error: no filters created")
	}

	return nil
}

// Returns true if address is in list of addresses to
// filter events for or if no filtering is specified
func (c *Contract) IsEventAddr(addr string) bool {
	if c.EventAddrs == nil {
		return false
	} else if len(c.EventAddrs) == 0 {
		return true
	} else if a, ok := c.EventAddrs[addr]; ok {
		return a
	}

	return false
}

// Returns true if address is in list of addresses to
// poll methods for or if no filtering is specified
func (c *Contract) IsMethodAddr(addr string) bool {
	if c.MethodAddrs == nil {
		return false
	} else if len(c.MethodAddrs) == 0 {
		return true
	} else if a, ok := c.MethodAddrs[addr]; ok {
		return a
	}

	return false
}

// Returns true if mapping value matches filtered for address or if not filter exists
// Used to check if an event log name-value mapping should be filtered or not
func (c *Contract) PassesEventFilter(args map[string]string) bool {
	for _, arg := range args {
		if c.IsEventAddr(arg) {
			return true
		}
	}

	return false
}

// Used to add an address to the token holder address list
// if it is on the method polling list or the filter is open
func (c *Contract) AddTokenHolderAddress(addr string) {
	if c.TknHolderAddrs != nil && c.IsMethodAddr(addr) {
		c.TknHolderAddrs[addr] = true
	}
}
