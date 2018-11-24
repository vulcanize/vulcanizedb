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

package parser

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

// Parser is used to fetch and parse contract ABIs
// It is dependent on etherscan's api
type Parser interface {
	Parse(contractAddr string) error
	Abi() string
	ParsedAbi() abi.ABI
	GetMethods(wanted []string) map[string]types.Method
	GetSelectMethods(wanted []string) map[string]types.Method
	GetEvents(wanted []string) map[string]types.Event
}

type parser struct {
	client    *geth.EtherScanAPI
	abi       string
	parsedAbi abi.ABI
}

func NewParser(network string) *parser {
	url := geth.GenURL(network)

	return &parser{
		client: geth.NewEtherScanClient(url),
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
func (p *parser) Parse(contractAddr string) error {
	knownAbi, err := p.lookUp(contractAddr)
	if err == nil {
		p.abi = knownAbi
		p.parsedAbi, err = geth.ParseAbi(knownAbi)
		return err
	}

	abiStr, err := p.client.GetAbi(contractAddr)
	if err != nil {
		return err
	}
	p.abi = abiStr
	p.parsedAbi, err = geth.ParseAbi(abiStr)

	return err
}

func (p *parser) lookUp(contractAddr string) (string, error) {
	if v, ok := constants.Abis[common.HexToAddress(contractAddr)]; ok {
		return v, nil
	}

	return "", errors.New("ABI not present in lookup tabe")
}

// Returns wanted methods, if they meet the criteria, as map of types.Methods
// Empty wanted array => all methods that fit are returned
// Nil wanted array => no events are returned
func (p *parser) GetSelectMethods(wanted []string) map[string]types.Method {
	addrMethods := map[string]types.Method{}
	if wanted == nil {
		return nil
	}

	for _, m := range p.parsedAbi.Methods {
		// Only return methods that have less than 3 inputs, 1 output
		if len(m.Inputs) < 3 && len(m.Outputs) == 1 && (len(wanted) == 0 || stringInSlice(wanted, m.Name)) {
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

// Returns wanted methods as map of types.Methods
// Empty wanted array => all events are returned
// Nil wanted array => no events are returned
func (p *parser) GetMethods(wanted []string) map[string]types.Method {
	methods := map[string]types.Method{}
	if wanted == nil {
		return methods
	}

	length := len(wanted)
	for _, m := range p.parsedAbi.Methods {
		if length == 0 || stringInSlice(wanted, m.Name) {
			methods[m.Name] = types.NewMethod(m)
		}
	}

	return methods
}

// Returns wanted events as map of types.Events
// Empty wanted array => all events are returned
// Nil wanted array => no events are returned
func (p *parser) GetEvents(wanted []string) map[string]types.Event {
	events := map[string]types.Event{}
	if wanted == nil {
		return events
	}

	length := len(wanted)
	for _, e := range p.parsedAbi.Events {
		if length == 0 || stringInSlice(wanted, e.Name) {
			events[e.Name] = types.NewEvent(e)
		}
	}

	return events
}

func wantType(arg abi.Argument) bool {
	wanted := []byte{abi.UintTy, abi.IntTy, abi.BoolTy, abi.StringTy, abi.AddressTy, abi.HashTy, abi.BytesTy, abi.FixedBytesTy, abi.FixedPointTy}
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
