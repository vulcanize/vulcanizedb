package geth

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func CalcUnclesReward(gethBlock *types.Block) float64 {
	var unclesReward float64
	for _, uncle := range gethBlock.Uncles() {
		blockNumber := gethBlock.Number().Int64()
		staticBlockReward := float64(staticRewardByBlockNumber(blockNumber))
		unclesReward += (1.0 + float64(uncle.Number.Int64()-gethBlock.Number().Int64())/8.0) * staticBlockReward
	}
	return unclesReward
}

func CalcBlockReward(gethBlock *types.Block, client GethClient) float64 {
	blockNumber := gethBlock.Number().Int64()
	staticBlockReward := staticRewardByBlockNumber(blockNumber)
	transactionFees := calcTransactionFees(gethBlock, client)
	uncleInclusionRewards := calcUncleInclusionRewards(gethBlock)
	return transactionFees + uncleInclusionRewards + staticBlockReward
}

func calcUncleInclusionRewards(gethBlock *types.Block) float64 {
	var uncleInclusionRewards float64
	staticBlockReward := staticRewardByBlockNumber(gethBlock.Number().Int64())
	for range gethBlock.Uncles() {
		uncleInclusionRewards += staticBlockReward * 1 / 32
	}
	return uncleInclusionRewards
}

func calcTransactionFees(gethBlock *types.Block, client GethClient) float64 {
	var transactionFees float64
	for _, transaction := range gethBlock.Transactions() {
		receipt, err := client.TransactionReceipt(context.Background(), transaction.Hash())
		if err != nil {
			continue
		}
		transactionFees += float64(transaction.GasPrice().Int64() * receipt.GasUsed.Int64())
	}
	return transactionFees / params.Ether
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
