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

//Our Block to DB
func BlockToBlockRecord(block Block) *BlockRecord {
	return &BlockRecord{
		BlockNumber: block.Number.Int64(),
		GasLimit:    block.GasLimit.Int64(),
		GasUsed:     block.GasUsed.Int64(),
		Time:        block.Time.Int64(),
	}
}

//DB block representation
type BlockRecord struct {
	BlockNumber int64 `db:"block_number"`
	GasLimit    int64 `db:"block_gaslimit"`
	GasUsed     int64 `db:"block_gasused"`
	Time        int64 `db:"block_time"`
}

//Geth Block to Ours
func GethBlockToCoreBlock(gethBlock *types.Block) Block {
	return Block{
		Number:               gethBlock.Number(),
		GasLimit:             gethBlock.GasLimit(),
		GasUsed:              gethBlock.GasUsed(),
		Time:                 gethBlock.Time(),
		NumberOfTransactions: len(gethBlock.Transactions()),
	}
}
