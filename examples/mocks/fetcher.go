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
	"errors"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Fetcher struct {
	ContractAddress string
	Abi             string
	FetchedBlocks   []int64
	BlockChain      core.BlockChain
	supply          big.Int
	balance         map[string]*big.Int
	allowance       map[string]map[string]*big.Int
	owner           common.Address
	stopped         bool
	stringName      string
	hashName        common.Hash
	stringSymbol    string
	hashSymbol      common.Hash
}

func (f *Fetcher) SetSupply(supply string) {
	f.supply.SetString(supply, 10)
}

func (f Fetcher) GetBlockChain() core.BlockChain {
	return f.BlockChain
}

func (f *Fetcher) FetchBigInt(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {

	f.Abi = contractAbi
	f.ContractAddress = contractAddress
	f.FetchedBlocks = append(f.FetchedBlocks, blockNumber)

	accumulator := big.NewInt(1)

	if method == "totalSupply" {
		f.supply.Add(&f.supply, accumulator)

		return f.supply, nil
	}

	if method == "balanceOf" {
		rfl := reflect.ValueOf(methodArgs[0])
		tokenHolderAddr := rfl.Interface().(string)
		pnt := f.balance[tokenHolderAddr]
		f.balance[tokenHolderAddr].Add(pnt, accumulator)

		return *f.balance[tokenHolderAddr], nil
	}

	if method == "allowance" {
		rfl1 := reflect.ValueOf(methodArgs[0])
		rfl2 := reflect.ValueOf(methodArgs[1])
		tokenHolderAddr := rfl1.Interface().(string)
		spenderAddr := rfl2.Interface().(string)
		pnt := f.allowance[tokenHolderAddr][spenderAddr]
		f.allowance[tokenHolderAddr][spenderAddr].Add(pnt, accumulator)

		return *f.allowance[tokenHolderAddr][spenderAddr], nil
	}

	return *big.NewInt(0), errors.New("invalid method argument")

}

func (f *Fetcher) FetchBool(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (bool, error) {

	f.Abi = contractAbi
	f.ContractAddress = contractAddress
	f.FetchedBlocks = append(f.FetchedBlocks, blockNumber)

	b := true

	if method == "stopped" {
		f.stopped = b

		return f.stopped, nil
	}

	return false, errors.New("invalid method argument")
}

func (f *Fetcher) FetchAddress(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Address, error) {

	f.Abi = contractAbi
	f.ContractAddress = contractAddress
	f.FetchedBlocks = append(f.FetchedBlocks, blockNumber)

	adr := common.HexToAddress("test_address")

	if method == "owner" {
		f.owner = adr

		return f.owner, nil
	}
	return common.HexToAddress(""), errors.New("invalid method argument")
}

func (f *Fetcher) FetchString(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (string, error) {

	f.Abi = contractAbi
	f.ContractAddress = contractAddress
	f.FetchedBlocks = append(f.FetchedBlocks, blockNumber)

	if method == "name" {
		f.stringName = "test_name"

		return f.stringName, nil
	}

	if method == "symbol" {
		f.stringSymbol = "test_symbol"

		return f.stringSymbol, nil
	}
	return "", errors.New("invalid method argument")
}

func (f *Fetcher) FetchHash(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Hash, error) {

	f.Abi = contractAbi
	f.ContractAddress = contractAddress
	f.FetchedBlocks = append(f.FetchedBlocks, blockNumber)

	if method == "name" {
		f.hashName = common.HexToHash("test_name")

		return f.hashName, nil
	}

	if method == "symbol" {
		f.hashSymbol = common.HexToHash("test_symbol")

		return f.hashSymbol, nil
	}
	return common.HexToHash(""), errors.New("invalid method argument")
}
