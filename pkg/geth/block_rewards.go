package geth

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
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
		transactionFees += float64(transaction.GasPrice * receipt.GasUsed)
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
