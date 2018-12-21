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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

type Poller interface {
	PollContract(con contract.Contract) error
	PollContractAt(con contract.Contract, blockNumber int64) error
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

func (p *poller) PollContract(con contract.Contract) error {
	for i := con.StartingBlock; i <= con.LastBlock; i++ {
		p.PollContractAt(con, i)
	}

	return nil
}

func (p *poller) PollContractAt(con contract.Contract, blockNumber int64) error {
	p.contract = con
	for _, m := range con.Methods {
		switch len(m.Args) {
		case 0:
			if err := p.pollNoArgAt(m, blockNumber); err != nil {
				return err
			}
		case 1:
			if err := p.pollSingleArgAt(m, blockNumber); err != nil {
				return err
			}
		case 2:
			if err := p.pollDoubleArgAt(m, blockNumber); err != nil {
				return err
			}
		default:
			return errors.New("poller error: too many arguments to handle")

		}
	}

	return nil
}

func (p *poller) pollNoArgAt(m types.Method, bn int64) error {
	result := types.Result{
		Block:  bn,
		Method: m,
		Inputs: nil,
		PgType: m.Return[0].PgType,
	}

	var out interface{}
	err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, nil, &out, bn)
	if err != nil {
		return errors.New(fmt.Sprintf("poller error calling 0 argument method\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", bn, m.Name, p.contract.Address, err))
	}

	strOut, err := stringify(out)
	if err != nil {
		return err
	}

	// Cache returned value if piping is turned on
	p.cache(out)
	result.Output = strOut

	// Persist result immediately
	err = p.PersistResults([]types.Result{result}, m, p.contract.Address, p.contract.Name)
	if err != nil {
		return errors.New(fmt.Sprintf("poller error persisting 0 argument method result\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", bn, m.Name, p.contract.Address, err))
	}

	return nil
}

// Use token holder address to poll methods that take 1 address argument (e.g. balanceOf)
func (p *poller) pollSingleArgAt(m types.Method, bn int64) error {
	result := types.Result{
		Block:  bn,
		Method: m,
		Inputs: make([]interface{}, 1),
		PgType: m.Return[0].PgType,
	}

	// Depending on the type of the arg choose
	// the correct argument set to iterate over
	var args map[interface{}]bool
	switch m.Args[0].Type.T {
	case abi.HashTy, abi.FixedBytesTy:
		args = p.contract.EmittedHashes
	case abi.AddressTy:
		args = p.contract.EmittedAddrs
	}
	if len(args) == 0 { // If we haven't collected any args by now we can't call the method
		return nil
	}
	results := make([]types.Result, 0, len(args))

	for arg := range args {
		in := []interface{}{arg}
		strIn := []interface{}{contract.StringifyArg(arg)}

		var out interface{}
		err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, in, &out, bn)
		if err != nil {
			return errors.New(fmt.Sprintf("poller error calling 1 argument method\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", bn, m.Name, p.contract.Address, err))
		}
		strOut, err := stringify(out)
		if err != nil {
			return err
		}
		p.cache(out)

		// Write inputs and outputs to result and append result to growing set
		result.Inputs = strIn
		result.Output = strOut
		results = append(results, result)
	}

	// Persist result set as batch
	err := p.PersistResults(results, m, p.contract.Address, p.contract.Name)
	if err != nil {
		return errors.New(fmt.Sprintf("poller error persisting 1 argument method result\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", bn, m.Name, p.contract.Address, err))
	}

	return nil
}

// Use token holder address to poll methods that take 2 address arguments (e.g. allowance)
func (p *poller) pollDoubleArgAt(m types.Method, bn int64) error {
	result := types.Result{
		Block:  bn,
		Method: m,
		Inputs: make([]interface{}, 2),
		PgType: m.Return[0].PgType,
	}

	// Depending on the type of the args choose
	// the correct argument sets to iterate over
	var firstArgs map[interface{}]bool
	switch m.Args[0].Type.T {
	case abi.HashTy, abi.FixedBytesTy:
		firstArgs = p.contract.EmittedHashes
	case abi.AddressTy:
		firstArgs = p.contract.EmittedAddrs
	}
	if len(firstArgs) == 0 {
		return nil
	}

	var secondArgs map[interface{}]bool
	switch m.Args[1].Type.T {
	case abi.HashTy, abi.FixedBytesTy:
		secondArgs = p.contract.EmittedHashes
	case abi.AddressTy:
		secondArgs = p.contract.EmittedAddrs
	}
	if len(secondArgs) == 0 {
		return nil
	}

	results := make([]types.Result, 0, len(firstArgs)*len(secondArgs))

	for arg1 := range firstArgs {
		for arg2 := range secondArgs {
			in := []interface{}{arg1, arg2}
			strIn := []interface{}{contract.StringifyArg(arg1), contract.StringifyArg(arg2)}

			var out interface{}
			err := p.bc.FetchContractData(p.contract.Abi, p.contract.Address, m.Name, in, &out, bn)
			if err != nil {
				return errors.New(fmt.Sprintf("poller error calling 2 argument method\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", bn, m.Name, p.contract.Address, err))
			}

			strOut, err := stringify(out)
			if err != nil {
				return err
			}

			p.cache(out)

			result.Output = strOut
			result.Inputs = strIn
			results = append(results, result)

		}
	}

	err := p.PersistResults(results, m, p.contract.Address, p.contract.Name)
	if err != nil {
		return errors.New(fmt.Sprintf("poller error persisting 2 argument method result\r\nblock: %d, method: %s, contract: %s\r\nerr: %v", bn, m.Name, p.contract.Address, err))
	}

	return nil
}

// This is just a wrapper around the poller blockchain's FetchContractData method
func (p *poller) FetchContractData(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	return p.bc.FetchContractData(contractAbi, contractAddress, method, methodArgs, result, blockNumber)
}

// This is used to cache an method return value if method piping is turned on
func (p *poller) cache(out interface{}) {
	// Cache returned value if piping is turned on
	if p.contract.Piping {
		switch out.(type) {
		case common.Hash:
			if p.contract.EmittedHashes != nil {
				p.contract.AddEmittedHash(out.(common.Hash))
			}
		case []byte:
			if p.contract.EmittedHashes != nil && len(out.([]byte)) == 32 {
				p.contract.AddEmittedHash(common.BytesToHash(out.([]byte)))
			}
		case common.Address:
			if p.contract.EmittedAddrs != nil {
				p.contract.AddEmittedAddr(out.(common.Address))
			}
		default:
		}
	}
}

func stringify(input interface{}) (string, error) {
	switch input.(type) {
	case *big.Int:
		b := input.(*big.Int)
		return b.String(), nil
	case common.Address:
		a := input.(common.Address)
		return a.String(), nil
	case common.Hash:
		h := input.(common.Hash)
		return h.String(), nil
	case string:
		return input.(string), nil
	case []byte:
		b := hexutil.Encode(input.([]byte))
		return b, nil
	case byte:
		b := input.(byte)
		return string(b), nil
	case bool:
		return strconv.FormatBool(input.(bool)), nil
	default:
		return "", errors.New("error: unhandled return type")
	}
}
