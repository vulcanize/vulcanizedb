// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package contract

import (
	"errors"
	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

type Contract struct {
	Name          string
	Address       string
	StartingBlock int64
	Abi           string
	ParsedAbi     abi.ABI
	Events        map[string]*types.Event      // Map of events to their names
	Methods       map[string]*types.Method     // Map of methods to their names
	Filters       map[string]filters.LogFilter // Map of event filters to their names
	Addresses     map[string]bool              // Map of all contract-associated addresses, populated as events are transformed
}

func (c *Contract) GenerateFilters(subset []string) error {
	c.Filters = map[string]filters.LogFilter{}
	for name, event := range c.Events {
		if len(subset) == 0 || stringInSlice(subset, name) {
			c.Filters[name] = filters.LogFilter{
				Name:      name,
				FromBlock: c.StartingBlock,
				ToBlock:   -1,
				Address:   c.Address,
				Topics:    core.Topics{helpers.GenerateSignature(event.Sig())},
			}
		}
	}

	if len(c.Filters) == 0 {
		return errors.New("error: no filters created")
	}

	return nil
}

func (c *Contract) AddAddress(addr string) {
	c.Addresses[addr] = true
}

func stringInSlice(list []string, s string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}

	return false
}
