// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// (U_n + 8 - B_n) * R / 8
// Returns a map of miner addresses to a map of the uncles they mined (hashes) to the rewards received for that uncle
func CalcUnclesReward(block core.Block, uncles []*types.Header) (*big.Int, map[string]map[string]*big.Int) {
	uncleRewards := new(big.Int)
	mappedUncleRewards := make(map[string]map[string]*big.Int)
	for _, uncle := range uncles {
		staticBlockReward := staticRewardByBlockNumber(block.Number)
		rewardDiv8 := staticBlockReward.Div(staticBlockReward, big.NewInt(8))
		uncleBlock := big.NewInt(uncle.Number.Int64())
		uncleBlockPlus8 := uncleBlock.Add(uncleBlock, big.NewInt(8))
		mainBlock := big.NewInt(block.Number)
		uncleBlockPlus8MinusMainBlock := uncleBlockPlus8.Sub(uncleBlockPlus8, mainBlock)
		thisUncleReward := rewardDiv8.Mul(rewardDiv8, uncleBlockPlus8MinusMainBlock)
		uncleRewards = uncleRewards.Add(uncleRewards, thisUncleReward)
		mappedUncleRewards[uncle.Coinbase.Hex()][uncle.Hash().Hex()].Add(mappedUncleRewards[uncle.Coinbase.Hex()][uncle.Hash().Hex()], thisUncleReward)
	}
	return uncleRewards, mappedUncleRewards
}

func CalcBlockReward(block core.Block, uncles []*types.Header) *big.Int {
	staticBlockReward := staticRewardByBlockNumber(block.Number)
	transactionFees := calcTransactionFees(block)
	uncleInclusionRewards := calcUncleInclusionRewards(block, uncles)
	tmp := transactionFees.Add(transactionFees, uncleInclusionRewards)
	return tmp.Add(tmp, staticBlockReward)
}

func calcTransactionFees(block core.Block) *big.Int {
	transactionFees := new(big.Int)
	for _, transaction := range block.Transactions {
		receipt := transaction.Receipt
		gasPrice := big.NewInt(transaction.GasPrice)
		gasUsed := big.NewInt(int64(receipt.GasUsed))
		transactionFee := gasPrice.Mul(gasPrice, gasUsed)
		transactionFees = transactionFees.Add(transactionFees, transactionFee)
	}
	return transactionFees
}

func calcUncleInclusionRewards(block core.Block, uncles []*types.Header) *big.Int {
	uncleInclusionRewards := new(big.Int)
	for range uncles {
		staticBlockReward := staticRewardByBlockNumber(block.Number)
		staticBlockReward.Div(staticBlockReward, big.NewInt(32))
		uncleInclusionRewards.Add(uncleInclusionRewards, staticBlockReward)
	}
	return uncleInclusionRewards
}

func staticRewardByBlockNumber(blockNumber int64) *big.Int {
	staticBlockReward := new(big.Int)
	//https://blog.ethereum.org/2017/10/12/byzantium-hf-announcement/
	if blockNumber >= 4370000 {
		staticBlockReward.SetString("3000000000000000000", 10)
	} else {
		staticBlockReward.SetString("5000000000000000000", 10)
	}
	return staticBlockReward
}
