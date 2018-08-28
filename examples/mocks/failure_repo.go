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
