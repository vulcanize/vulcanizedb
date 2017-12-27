package geth

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func UncleReward(gethBlock *types.Block, client GethClient) float64 {
	var uncleReward float64
	for _, uncle := range gethBlock.Uncles() {
		staticBlockReward := float64(blockNumberStaticReward(gethBlock)) / float64(8)
		uncleReward += float64(uncle.Number.Int64()-gethBlock.Number().Int64()+int64(8)) * staticBlockReward
	}
	return uncleReward
}

func BlockReward(gethBlock *types.Block, client GethClient) float64 {
	staticBlockReward := blockNumberStaticReward(gethBlock)
	transactionFees := calcTransactionFees(gethBlock, client)
	uncleInclusionRewards := uncleInclusionRewards(gethBlock, staticBlockReward)
	return transactionFees + uncleInclusionRewards + staticBlockReward
}

func uncleInclusionRewards(gethBlock *types.Block, staticBlockReward float64) float64 {
	var uncleInclusionRewards float64
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

func blockNumberStaticReward(gethBlock *types.Block) float64 {
	var staticBlockReward float64
	if gethBlock.Number().Int64() > 4269999 {
		staticBlockReward = 3
	} else {
		staticBlockReward = 5
	}
	return staticBlockReward
}
