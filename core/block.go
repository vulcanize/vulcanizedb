package core

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number       int64
	GasLimit     int64
	GasUsed      int64
	Time         int64
	Transactions []Transaction
}

func GethBlockToCoreBlock(gethBlock *types.Block) Block {
	transactions := []Transaction{}
	for _, gethTransaction := range gethBlock.Transactions() {
		transactions = append(transactions, gethTransToCoreTrans(gethTransaction))
	}
	return Block{
		Number:       gethBlock.Number().Int64(),
		GasLimit:     gethBlock.GasLimit().Int64(),
		GasUsed:      gethBlock.GasUsed().Int64(),
		Time:         gethBlock.Time().Int64(),
		Transactions: transactions,
	}
}
