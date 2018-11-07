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
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

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
