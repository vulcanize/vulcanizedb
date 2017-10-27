package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

//Our block representation
type Block struct {
	Number               *big.Int
	GasLimit             *big.Int
	GasUsed              *big.Int
	Time                 *big.Int
	NumberOfTransactions int
}

//Geth Block to Ours
func GethBlockToCoreBlock(gethBlock *types.Block) Block {
	return Block{
		Number:               gethBlock.Number(),
		GasLimit:             gethBlock.GasLimit(),
		GasUsed:              gethBlock.GasUsed(),
		Time:                 gethBlock.Time(),
		NumberOfTransactions: gethBlock.Transactions().Len(),
	}
}
