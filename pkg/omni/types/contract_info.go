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

package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type ContractInfo struct {
	Name          string
	Address       string
	StartingBlock int64
	Abi           string
	ParsedAbi     abi.ABI
	Events        map[string]*Event            // Map of events to their names
	Methods       map[string]*Method           // Map of methods to their names
	Filters       map[string]filters.LogFilter // Map of event filters to their names
}

func (i *ContractInfo) GenerateFilters(subset []string) {
	i.Filters = map[string]filters.LogFilter{}
	for name, event := range i.Events {
		if len(subset) == 0 || stringInSlice(subset, name) {
			i.Filters[name] = filters.LogFilter{
				Name:      name,
				FromBlock: i.StartingBlock,
				ToBlock:   -1,
				Address:   i.Address,
				Topics:    core.Topics{event.Sig()},
			}
		}
	}
}

func stringInSlice(list []string, s string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}

	return false
}
