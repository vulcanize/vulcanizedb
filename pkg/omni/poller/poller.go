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

package poller

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/omni/constants"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

type Poller interface {
	PollContract(con *contract.Contract) error
	PollMethod(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error
}

type poller struct {
	bc       core.BlockChain
	contract *contract.Contract
}

func NewPoller(blockChain core.BlockChain) *poller {

	return &poller{
		bc: blockChain,
	}
}

// Used to call contract's methods found in abi using list of contract-related addresses
func (p *poller) PollContract(con *contract.Contract) error {
	p.contract = con
	// Iterate over each of the contracts methods
	for _, m := range con.Methods {
		switch len(m.Inputs) {
		case 0:
			if err := p.pollNoArg(m); err != nil {
				return err
			}
		case 1:
			if err := p.pollSingleArg(m); err != nil {
				return err
			}
		case 2:
			if err := p.pollDoubleArg(m); err != nil {
				return err
			}
		default:
			return errors.New("poller error: too many arguments to handle")

		}
	}

	return nil
}

// Poll methods that take no arguments
func (p *poller) pollNoArg(m *types.Method) error {
	result := &types.Result{
		Inputs:  nil,
		Outputs: map[int64]interface{}{},
		PgType:  m.Outputs[0].PgType,
	}

	for i := p.contract.StartingBlock; i <= p.contract.LastBlock; i++ {
		var res interface{}
		err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, result.Inputs, &res, i)
		if err != nil {
			return err
		}
		result.Outputs[i] = res
	}

	// Persist results now instead of holding onto them
	m.Results = append(m.Results, result)

	return nil
}

// Use token holder address to poll methods that take 1 address argument (e.g. balanceOf)
func (p *poller) pollSingleArg(m *types.Method) error {
	for addr := range p.contract.TknHolderAddrs {
		result := &types.Result{
			Inputs:  make([]interface{}, 1),
			Outputs: map[int64]interface{}{},
			PgType:  m.Outputs[0].PgType,
		}
		result.Inputs[0] = common.HexToAddress(addr)

		for i := p.contract.StartingBlock; i <= p.contract.LastBlock; i++ {
			var res interface{}
			err := p.bc.FetchContractData(constants.TusdAbiString, p.contract.Address, m.Name, result.Inputs, &res, i)
			if err != nil {
				return err
			}
			result.Outputs[i] = res
		}

		m.Results = append(m.Results, result)
	}

	return nil
}

// Use token holder address to poll methods that take 2 address arguments (e.g. allowance)
func (p *poller) pollDoubleArg(m *types.Method) error {
	// For a large block range and address list this will take a really, really long time- maybe we should only do 1 arg methods
	for addr1 := range p.contract.TknHolderAddrs {
		for addr2 := range p.contract.TknHolderAddrs {
			result := &types.Result{
				Inputs:  make([]interface{}, 2),
				Outputs: map[int64]interface{}{},
				PgType:  m.Outputs[0].PgType,
			}
			result.Inputs[0] = common.HexToAddress(addr1)
			result.Inputs[1] = common.HexToAddress(addr2)

			for i := p.contract.StartingBlock; i <= p.contract.LastBlock; i++ {
				var res interface{}
				err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, result.Inputs, &res, i)
				if err != nil {
					return err
				}
				result.Outputs[i] = res
			}

			m.Results = append(m.Results, result)
		}
	}

	return nil
}

// This is just a wrapper around the poller blockchain's FetchContractData method
func (p *poller) PollMethod(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	return p.bc.FetchContractData(contractAbi, contractAddress, method, methodArgs, result, blockNumber)
}
