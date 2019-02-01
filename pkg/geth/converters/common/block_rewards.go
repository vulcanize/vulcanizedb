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

package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func CalcUnclesReward(block core.Block, uncles []*types.Header) float64 {
	var unclesReward float64
	for _, uncle := range uncles {
		blockNumber := block.Number
		staticBlockReward := float64(staticRewardByBlockNumber(blockNumber))
		unclesReward += (1.0 + float64(uncle.Number.Int64()-block.Number)/8.0) * staticBlockReward
	}
	return unclesReward
}

func CalcBlockReward(block core.Block, uncles []*types.Header) float64 {
	blockNumber := block.Number
	staticBlockReward := staticRewardByBlockNumber(blockNumber)
	transactionFees := calcTransactionFees(block)
	uncleInclusionRewards := calcUncleInclusionRewards(block, uncles)
	return transactionFees + uncleInclusionRewards + staticBlockReward
}

func calcTransactionFees(block core.Block) float64 {
	var transactionFees float64
	for _, transaction := range block.Transactions {
		receipt := transaction.Receipt
		transactionFees += float64(uint64(transaction.GasPrice) * receipt.GasUsed)
	}
	return transactionFees / params.Ether
}

func calcUncleInclusionRewards(block core.Block, uncles []*types.Header) float64 {
	var uncleInclusionRewards float64
	staticBlockReward := staticRewardByBlockNumber(block.Number)
	for range uncles {
		uncleInclusionRewards += staticBlockReward * 1 / 32
	}
	return uncleInclusionRewards
}

func staticRewardByBlockNumber(blockNumber int64) float64 {
	var staticBlockReward float64
	//https://blog.ethereum.org/2017/10/12/byzantium-hf-announcement/
	if blockNumber >= 4370000 {
		staticBlockReward = 3
	} else {
		staticBlockReward = 5
	}
	return staticBlockReward
}
