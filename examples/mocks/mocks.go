// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocks

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"reflect"
)

type Fetcher struct {
	ContractAddress string
	Abi             string
	FetchedBlocks   []int64
	BlockChain      core.BlockChain
	supply          big.Int
	balance         map[string]*big.Int
	allowance       map[string]map[string]*big.Int
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
		rfl := reflect.ValueOf(methodArgs).Field(0)
		tokenHolderAddr := rfl.Interface().(string)
		pnt := f.balance[tokenHolderAddr]
		f.balance[tokenHolderAddr].Add(pnt, accumulator)

		return *f.balance[tokenHolderAddr], nil
	}

	if method == "allowance" {
		rfl1 := reflect.ValueOf(methodArgs).Field(0)
		rfl2 := reflect.ValueOf(methodArgs).Field(1)
		tokenHolderAddr := rfl1.Interface().(string)
		spenderAddr := rfl2.Interface().(string)
		pnt := f.allowance[tokenHolderAddr][spenderAddr]
		f.allowance[tokenHolderAddr][spenderAddr].Add(pnt, accumulator)

		return *f.allowance[tokenHolderAddr][spenderAddr], nil
	}

	return *big.NewInt(0), errors.New("invalid method argument")

}

func (f *Fetcher) FetchBool(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (bool, error) {
	//TODO: this
	return false, nil
}

func (f *Fetcher) FetchAddress(method, contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (common.Address, error) {
	//TODO: this
	var adr common.Address
	return adr, nil
}

type ERC20TokenRepository struct {
	TotalSuppliesCreated         []every_block.TokenSupply
	MissingSupplyBlockNumbers    []int64
	TotalBalancesCreated         []every_block.TokenBalance
	MissingBalanceBlockNumbers   []int64
	TotalAllowancesCreated       []every_block.TokenAllowance
	MissingAllowanceBlockNumbers []int64
	StartingBlock                int64
	EndingBlock                  int64
}

func (fr *ERC20TokenRepository) CreateSupply(supply every_block.TokenSupply) error {
	fr.TotalSuppliesCreated = append(fr.TotalSuppliesCreated, supply)
	return nil
}

func (fr *ERC20TokenRepository) CreateBalance(balance every_block.TokenBalance) error {
	fr.TotalBalancesCreated = append(fr.TotalBalancesCreated, balance)
	return nil
}

func (fr *ERC20TokenRepository) CreateAllowance(allowance every_block.TokenAllowance) error {
	fr.TotalAllowancesCreated = append(fr.TotalAllowancesCreated, allowance)
	return nil
}

func (fr *ERC20TokenRepository) MissingSupplyBlocks(startingBlock, highestBlock int64, tokenAddress string) ([]int64, error) {
	fr.StartingBlock = startingBlock
	fr.EndingBlock = highestBlock
	return fr.MissingSupplyBlockNumbers, nil
}

func (fr *ERC20TokenRepository) MissingBalanceBlocks(startingBlock, highestBlock int64, tokenAddress, holderAddress string) ([]int64, error) {
	fr.StartingBlock = startingBlock
	fr.EndingBlock = highestBlock
	return fr.MissingBalanceBlockNumbers, nil
}

func (fr *ERC20TokenRepository) MissingAllowanceBlocks(startingBlock, highestBlock int64, tokenAddress, holderAddress, spenderAddress string) ([]int64, error) {
	fr.StartingBlock = startingBlock
	fr.EndingBlock = highestBlock
	return fr.MissingAllowanceBlockNumbers, nil
}

func (fr *ERC20TokenRepository) SetMissingSupplyBlocks(missingBlocks []int64) {
	fr.MissingSupplyBlockNumbers = missingBlocks
}

func (fr *ERC20TokenRepository) SetMissingBalanceBlocks(missingBlocks []int64) {
	fr.MissingBalanceBlockNumbers = missingBlocks
}

func (fr *ERC20TokenRepository) SetMissingAllowanceBlocks(missingBlocks []int64) {
	fr.MissingAllowanceBlockNumbers = missingBlocks
}

type FailureRepository struct {
	createSupplyFail              bool
	createBalanceFail             bool
	createAllowanceFail           bool
	missingSupplyBlocksFail       bool
	missingBalanceBlocksFail      bool
	missingAllowanceBlocksFail    bool
	missingSupplyBlocksNumbers    []int64
	missingBalanceBlocksNumbers   []int64
	missingAllowanceBlocksNumbers []int64
}

func (fr *FailureRepository) CreateSupply(supply every_block.TokenSupply) error {
	if fr.createSupplyFail {
		return fakes.FakeError
	} else {
		return nil
	}
}

func (fr *FailureRepository) CreateBalance(balance every_block.TokenBalance) error {
	if fr.createBalanceFail {
		return fakes.FakeError
	} else {
		return nil
	}
}

func (fr *FailureRepository) CreateAllowance(allowance every_block.TokenAllowance) error {
	if fr.createAllowanceFail {
		return fakes.FakeError
	} else {
		return nil
	}
}

func (fr *FailureRepository) MissingSupplyBlocks(startingBlock, highestBlock int64, tokenAddress string) ([]int64, error) {
	if fr.missingSupplyBlocksFail {
		return []int64{}, fakes.FakeError
	} else {
		return fr.missingSupplyBlocksNumbers, nil
	}
}

func (fr *FailureRepository) MissingBalanceBlocks(startingBlock, highestBlock int64, tokenAddress, holderAddress string) ([]int64, error) {
	if fr.missingBalanceBlocksFail {
		return []int64{}, fakes.FakeError
	} else {
		return fr.missingBalanceBlocksNumbers, nil
	}
}

func (fr *FailureRepository) MissingAllowanceBlocks(startingBlock, highestBlock int64, tokenAddress, holderAddress, spenderAddress string) ([]int64, error) {
	if fr.missingAllowanceBlocksFail {
		return []int64{}, fakes.FakeError
	} else {
		return fr.missingAllowanceBlocksNumbers, nil
	}
}

func (fr *FailureRepository) SetCreateSupplyFail(fail bool) {
	fr.createSupplyFail = fail
}

func (fr *FailureRepository) SetCreateBalanceFail(fail bool) {
	fr.createBalanceFail = fail
}

func (fr *FailureRepository) SetCreateAllowanceFail(fail bool) {
	fr.createAllowanceFail = fail
}

func (fr *FailureRepository) SetMissingSupplyBlocksFail(fail bool) {
	fr.missingSupplyBlocksFail = fail
}

func (fr *FailureRepository) SetMissingBalanceBlocksFail(fail bool) {
	fr.missingBalanceBlocksFail = fail
}

func (fr *FailureRepository) SetMissingAllowanceBlocksFail(fail bool) {
	fr.missingAllowanceBlocksFail = fail
}

func (fr *FailureRepository) SetMissingSupplyBlocks(missingBlocks []int64) {
	fr.missingSupplyBlocksNumbers = missingBlocks
}

func (fr *FailureRepository) SetMissingBalanceBlocks(missingBlocks []int64) {
	fr.missingBalanceBlocksNumbers = missingBlocks
}

func (fr *FailureRepository) SetMissingAllowanceBlocks(missingBlocks []int64) {
	fr.missingAllowanceBlocksNumbers = missingBlocks
}
