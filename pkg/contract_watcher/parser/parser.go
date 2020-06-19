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

package parser

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/types"
	"github.com/makerdao/vulcanizedb/pkg/eth"
)

// Parser is used to fetch and parse contract ABIs
// It is dependent on etherscan's api
type Parser interface {
	Parse(contractAddr string) error
	ParseAbiStr(abiStr string) error
	Abi() string
	ParsedAbi() abi.ABI
	GetEvents(wanted []string) map[string]types.Event
}

type parser struct {
	client    *eth.EtherScanAPI
	abi       string
	parsedAbi abi.ABI
}

// NewParser returns a new Parser
func NewParser(network string) Parser {
	url := eth.GenURL(network)

	return &parser{
		client: eth.NewEtherScanClient(url),
	}
}

// Abi returns the parser's configured abi string
func (p *parser) Abi() string {
	return p.abi
}

// ParsedAbi returns the parser's parsed abi
func (p *parser) ParsedAbi() abi.ABI {
	return p.parsedAbi
}

// Parse retrieves and parses the abi string
// for the given contract address
func (p *parser) Parse(contractAddr string) error {
	// If the abi is one our locally stored abis, fetch
	// TODO: Allow users to pass abis through config
	knownAbi, err := p.lookUp(contractAddr)
	if err == nil {
		p.abi = knownAbi
		p.parsedAbi, err = eth.ParseAbi(knownAbi)
		return err
	}
	// Try getting abi from etherscan
	abiStr, err := p.client.GetAbi(contractAddr)
	if err != nil {
		return err
	}
	//TODO: Implement other ways to fetch abi
	p.abi = abiStr
	p.parsedAbi, err = eth.ParseAbi(abiStr)

	return err
}

// ParseAbiStr loads and parses an abi from a given abi string
func (p *parser) ParseAbiStr(abiStr string) error {
	var err error
	p.abi = abiStr
	p.parsedAbi, err = eth.ParseAbi(abiStr)

	return err
}

func (p *parser) lookUp(contractAddr string) (string, error) {
	if v, ok := constants.ABIs[common.HexToAddress(contractAddr)]; ok {
		return v, nil
	}

	return "", errors.New("ABI not present in lookup table")
}

// GetEvents returns wanted events as map of types.Events
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

func stringInSlice(list []string, s string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}

	return false
}
