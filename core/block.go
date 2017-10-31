package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Number       *big.Int
	GasLimit     *big.Int
	GasUsed      *big.Int
	Time         *big.Int
	Transactions []Transaction
}

func GethBlockToCoreBlock(gethBlock *types.Block) Block {
	transactions := []Transaction{}
	for _, gethTransaction := range gethBlock.Transactions() {
		transactions = append(transactions, gethTransToCoreTrans(gethTransaction))
	}
	return Block{
		Number:       gethBlock.Number(),
		GasLimit:     gethBlock.GasLimit(),
		GasUsed:      gethBlock.GasUsed(),
		Time:         gethBlock.Time(),
		Transactions: transactions,
	}
}
