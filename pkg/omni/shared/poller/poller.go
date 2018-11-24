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
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

type Poller interface {
	PollContract(con contract.Contract) error
	FetchContractData(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error
}

type poller struct {
	repository.MethodRepository
	bc       core.BlockChain
	contract contract.Contract
}

func NewPoller(blockChain core.BlockChain, db *postgres.DB, mode types.Mode) *poller {

	return &poller{
		MethodRepository: repository.NewMethodRepository(db, mode),
		bc:               blockChain,
	}
}

// Used to call contract's methods found in abi using list of contract-related addresses
func (p *poller) PollContract(con contract.Contract) error {
	p.contract = con
	// Iterate over each of the contracts methods
	for _, m := range con.Methods {
		switch len(m.Args) {
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
func (p *poller) pollNoArg(m types.Method) error {
	result := types.Result{
		Method: m,
		Inputs: nil,
		PgType: m.Return[0].PgType,
	}

	for i := p.contract.StartingBlock; i <= p.contract.LastBlock; i++ {
		var out interface{}
		err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, nil, &out, i)
		if err != nil {
			return errors.New(fmt.Sprintf("poller error calling 0 argument method\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", i, m.Name, p.contract.Address, err))
		}

		strOut, err := stringify(out)
		if err != nil {
			return err
		}

		result.Output = strOut
		result.Block = i

		// Persist result immediately
		err = p.PersistResult(result, p.contract.Address, p.contract.Name)
		if err != nil {
			return errors.New(fmt.Sprintf("poller error persisting 0 argument method result\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", i, m.Name, p.contract.Address, err))
		}
	}

	return nil
}

// Use token holder address to poll methods that take 1 address argument (e.g. balanceOf)
func (p *poller) pollSingleArg(m types.Method) error {
	result := types.Result{
		Method: m,
		Inputs: make([]interface{}, 1),
		PgType: m.Return[0].PgType,
	}

	for addr := range p.contract.TknHolderAddrs {
		for i := p.contract.StartingBlock; i <= p.contract.LastBlock; i++ {
			hashArgs := []common.Address{common.HexToAddress(addr)}
			in := make([]interface{}, len(hashArgs))
			strIn := make([]interface{}, len(hashArgs))
			for i, s := range hashArgs {
				in[i] = s
				strIn[i] = s.String()
			}

			var out interface{}
			err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, in, &out, i)
			if err != nil {
				return errors.New(fmt.Sprintf("poller error calling 1 argument method\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", i, m.Name, p.contract.Address, err))
			}

			strOut, err := stringify(out)
			if err != nil {
				return err
			}

			result.Output = strOut
			result.Block = i
			result.Inputs = strIn

			err = p.PersistResult(result, p.contract.Address, p.contract.Name)
			if err != nil {
				return errors.New(fmt.Sprintf("poller error persisting 1 argument method result\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", i, m.Name, p.contract.Address, err))
			}
		}
	}

	return nil
}

// Use token holder address to poll methods that take 2 address arguments (e.g. allowance)
func (p *poller) pollDoubleArg(m types.Method) error {
	// For a large block range and address list this will take a really, really long time- maybe we should only do 1 arg methods
	result := types.Result{
		Method: m,
		Inputs: make([]interface{}, 2),
		PgType: m.Return[0].PgType,
	}

	for addr1 := range p.contract.TknHolderAddrs {
		for addr2 := range p.contract.TknHolderAddrs {
			for i := p.contract.StartingBlock; i <= p.contract.LastBlock; i++ {
				hashArgs := []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}
				in := make([]interface{}, len(hashArgs))
				strIn := make([]interface{}, len(hashArgs))
				for i, s := range hashArgs {
					in[i] = s
					strIn[i] = s.String()
				}

				var out interface{}
				err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, in, &out, i)
				if err != nil {
					return errors.New(fmt.Sprintf("poller error calling 2 argument method\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", i, m.Name, p.contract.Address, err))
				}

				strOut, err := stringify(out)
				if err != nil {
					return err
				}

				result.Output = strOut
				result.Block = i
				result.Inputs = strIn

				err = p.PersistResult(result, p.contract.Address, p.contract.Name)
				if err != nil {
					return errors.New(fmt.Sprintf("poller error persisting 2 argument method result\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", i, m.Name, p.contract.Address, err))
				}
			}

		}
	}

	return nil
}

// This is just a wrapper around the poller blockchain's FetchContractData method
func (p *poller) FetchContractData(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	return p.bc.FetchContractData(contractAbi, contractAddress, method, methodArgs, result, blockNumber)
}

func stringify(input interface{}) (string, error) {
	switch input.(type) {
	case *big.Int:
		var b *big.Int
		b = input.(*big.Int)
		return b.String(), nil
	case common.Address:
		var a common.Address
		a = input.(common.Address)
		return a.String(), nil
	case common.Hash:
		var h common.Hash
		h = input.(common.Hash)
		return h.String(), nil
	case string:
		return input.(string), nil
	case bool:
		return strconv.FormatBool(input.(bool)), nil
	default:
		return "", errors.New("error: unhandled return type")
	}
}
