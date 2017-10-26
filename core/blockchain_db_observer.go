package core

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BlockchainDBObserver struct {
	Db *sqlx.DB
}

func (observer BlockchainDBObserver) NotifyBlockAdded(block Block) {
	blockRecord := BlockToBlockRecord(block)
	observer.Db.NamedExec(
		"INSERT INTO blocks "+
			"(block_number, block_gaslimit, block_gasused, block_time) "+
			"VALUES (:block_number, :block_gaslimit, :block_gasused, :block_time)", blockRecord)
	//observer.Db.MustExec("Insert INTO blocks (block_number) VALUES ($1)", block.Number.Int64())
}
