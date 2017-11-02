package observers

import (
	"github.com/8thlight/vulcanizedb/core"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BlockchainDBObserver struct {
	Db *sqlx.DB
}

func (observer BlockchainDBObserver) NotifyBlockAdded(block core.Block) {
	insertedBlockId := saveBlock(observer, block)
	saveTransactions(insertedBlockId, block.Transactions, observer)
}

func saveBlock(observer BlockchainDBObserver, block core.Block) int64 {
	insertedBlock := observer.Db.QueryRow(
		"Insert INTO blocks "+
		"(block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		block.Number, block.GasLimit, block.GasUsed, block.Time, block.Difficulty, block.Hash, block.Nonce, block.ParentHash, block.Size, block.UncleHash)
	var blockId int64
	insertedBlock.Scan(&blockId)
	return blockId
}

func saveTransactions(blockId int64, transactions []core.Transaction, observer BlockchainDBObserver) {
	for _, transaction := range transactions {
		observer.Db.MustExec("Insert INTO transactions "+
			"(block_id, tx_hash, tx_nonce, tx_to, tx_gaslimit, tx_gasprice, tx_value) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			blockId, transaction.Hash, transaction.Nonce, transaction.To, transaction.GasLimit, transaction.GasPrice, transaction.Value)
	}
}
