package core

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BlockchainDBObserver struct {
	Db *sqlx.DB
}

func (observer BlockchainDBObserver) NotifyBlockAdded(block Block) {
	observer.Db.MustExec("Insert INTO blocks "+
		"(block_number, block_gaslimit, block_gasused, block_time) "+
		"VALUES ($1, $2, $3, $4)",
		block.Number.Int64(), block.GasLimit.Int64(), block.GasUsed.Int64(), block.Time.Int64())
}
