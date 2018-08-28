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
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
)

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
