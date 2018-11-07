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
