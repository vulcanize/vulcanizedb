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

package parser

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

// Parser is used to fetch and parse contract ABIs
// It is dependent on etherscan's api
type Parser interface {
	Parse(contractAddr string) error
	Abi() string
	ParsedAbi() abi.ABI
	GetMethods() map[string]*types.Method
	GetEvents() map[string]*types.Event
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
	abiStr, err := p.client.GetAbi(contractAddr)
	if err != nil {
		return err
	}

	p.abi = abiStr
	p.parsedAbi, err = geth.ParseAbi(abiStr)

	return err
}

// Parses methods into our custom method type and returns
func (p *parser) GetMethods() map[string]*types.Method {
	methods := map[string]*types.Method{}

	for _, m := range p.parsedAbi.Methods {
		method := types.NewMethod(m)
		methods[m.Name] = method
	}

	return methods
}

// Parses events into our custom event type and returns
func (p *parser) GetEvents() map[string]*types.Event {
	events := map[string]*types.Event{}

	for _, e := range p.parsedAbi.Events {
		event := types.NewEvent(e)
		events[e.Name] = event
	}

	return events
}
