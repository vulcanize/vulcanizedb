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
	"math/big"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

type Fetcher struct {
	ContractAddress string
	Abi             string
	FetchedBlocks   []int64
	BlockChain      core.BlockChain
	supply          big.Int
}

func (f *Fetcher) SetSupply(supply string) {
	f.supply.SetString(supply, 10)
}

func (f Fetcher) GetBlockChain() core.BlockChain {
	return f.BlockChain
}

func (f *Fetcher) FetchSupplyOf(contractAbi string, contractAddress string, blockNumber int64) (big.Int, error) {
	f.Abi = contractAbi
	f.ContractAddress = contractAddress
	f.FetchedBlocks = append(f.FetchedBlocks, blockNumber)

	accumulator := big.NewInt(1)
	f.supply.Add(&f.supply, accumulator)

	return f.supply, nil
}

type TotalSupplyRepository struct {
	TotalSuppliesCreated []every_block.TokenSupply
	MissingBlockNumbers  []int64
	StartingBlock        int64
	EndingBlock          int64
}

func (fr *TotalSupplyRepository) Create(supply every_block.TokenSupply) error {
	fr.TotalSuppliesCreated = append(fr.TotalSuppliesCreated, supply)
	return nil
}

func (fr *TotalSupplyRepository) MissingBlocks(startingBlock int64, highestBlock int64) ([]int64, error) {
	fr.StartingBlock = startingBlock
	fr.EndingBlock = highestBlock
	return fr.MissingBlockNumbers, nil
}

func (fr *TotalSupplyRepository) SetMissingBlocks(missingBlocks []int64) {
	fr.MissingBlockNumbers = missingBlocks
}

type FailureRepository struct {
	createFail           bool
	missingBlocksFail    bool
	missingBlocksNumbers []int64
}

func (fr *FailureRepository) Create(supply every_block.TokenSupply) error {
	if fr.createFail {
		return fakes.FakeError
	} else {
		return nil
	}
}

func (fr *FailureRepository) MissingBlocks(startingBlock int64, highestBlock int64) ([]int64, error) {
	if fr.missingBlocksFail {
		return []int64{}, fakes.FakeError
	} else {
		return fr.missingBlocksNumbers, nil
	}
}

func (fr *FailureRepository) SetCreateFail(fail bool) {
	fr.createFail = fail
}

func (fr *FailureRepository) SetMissingBlocksFail(fail bool) {
	fr.missingBlocksFail = fail
}

func (fr *FailureRepository) SetMissingBlocks(missingBlocks []int64) {
	fr.missingBlocksNumbers = missingBlocks
}
