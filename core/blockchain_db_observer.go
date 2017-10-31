package core

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BlockchainDBObserver struct {
	Db *sqlx.DB
}

func (observer BlockchainDBObserver) NotifyBlockAdded(block Block) {
	insertedBlockId := saveBlock(observer, block)
	saveTransactions(insertedBlockId, block.Transactions, observer)
}

func saveBlock(observer BlockchainDBObserver, block Block) int64 {
	insertedBlock := observer.Db.QueryRow("Insert INTO blocks "+
		"(block_number, block_gaslimit, block_gasused, block_time) "+
		"VALUES ($1, $2, $3, $4) RETURNING id",
		block.Number.Int64(), block.GasLimit.Int64(), block.GasUsed.Int64(), block.Time.Int64())
	var blockId int64
	insertedBlock.Scan(&blockId)
	return blockId
}

func saveTransactions(blockId int64, transactions []Transaction, observer BlockchainDBObserver) {
	for _, transaction := range transactions {
		observer.Db.MustExec("Insert INTO transactions "+
			"(block_id, tx_hash, tx_nonce, tx_to, tx_gaslimit, tx_gasprice, tx_value) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			blockId, transaction.Hash, transaction.Nonce, transaction.To, transaction.GasLimit, transaction.GasPrice, transaction.Value)
	}
}
