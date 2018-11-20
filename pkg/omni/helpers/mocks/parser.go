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

package mocks

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
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
	p.parsedAbi, err = geth.ParseAbi(p.abi)

	return err
}

// Returns wanted methods, if they meet the criteria, as map of types.Methods
// Only returns specified methods
func (p *parser) GetMethods(wanted []string) map[string]types.Method {
	addrMethods := map[string]types.Method{}

	for _, m := range p.parsedAbi.Methods {
		// Only return methods that have less than 3 inputs, 1 output, and wanted
		if len(m.Inputs) < 3 && len(m.Outputs) == 1 && stringInSlice(wanted, m.Name) {
			addrsOnly := true
			for _, input := range m.Inputs {
				if input.Type.T != abi.AddressTy {
					addrsOnly = false
				}
			}

			// Only return methods if inputs are all of type address and output is of the accepted types
			if addrsOnly && wantType(m.Outputs[0]) {
				method := types.NewMethod(m)
				addrMethods[method.Name] = method
			}
		}
	}

	return addrMethods
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

func wantType(arg abi.Argument) bool {
	wanted := []byte{abi.UintTy, abi.IntTy, abi.BoolTy, abi.StringTy, abi.AddressTy, abi.HashTy}
	for _, ty := range wanted {
		if arg.Type.T == ty {
			return true
		}
	}

	return false
}

func stringInSlice(list []string, s string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}

	return false
}
