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

package mocks

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/types"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

// Mock parser
// Is given ABI string instead of address
// Performs all other functions of the real parser
type parser struct {
	abi       string
	parsedAbi abi.ABI
}

func NewParser(abi string) *parser {

	return &parser{
		abi: abi,
	}
}

func (p *parser) Abi() string {
	return p.abi
}

func (p *parser) ParsedAbi() abi.ABI {
	return p.parsedAbi
}

// Retrieves and parses the abi string
// for the given contract address
func (p *parser) Parse() error {
	var err error
	p.parsedAbi, err = eth.ParseAbi(p.abi)

	return err
}

// Returns wanted events as map of types.Events
// If no events are specified, all events are returned
func (p *parser) GetEvents(wanted []string) map[string]types.Event {
	events := map[string]types.Event{}

	for _, e := range p.parsedAbi.Events {
		if len(wanted) == 0 || stringInSlice(wanted, e.Name) {
			event := types.NewEvent(e)
			events[e.Name] = event
		}
	}

	return events
}

func stringInSlice(list []string, s string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}

	return false
}
