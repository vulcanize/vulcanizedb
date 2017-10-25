package core

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BlockRecord struct {
	BlockNumber int64 `db:"block_number"`
}

type BlockchainDBObserver struct {
	Db *sqlx.DB
}

func (observer BlockchainDBObserver) NotifyBlockAdded(block Block) {
	observer.Db.NamedExec("INSERT INTO blocks (block_number) VALUES (:block_number)", &BlockRecord{BlockNumber: block.Number.Int64()})
	//observer.Db.MustExec("Insert INTO blocks (block_number) VALUES ($1)", block.Number.Int64())
}
