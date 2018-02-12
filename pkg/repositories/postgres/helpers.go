package postgres

import (
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func (db *DB) clearData() {
	db.MustExec("DELETE FROM watched_contracts")
	db.MustExec("DELETE FROM transactions")
	db.MustExec("DELETE FROM blocks")
	db.MustExec("DELETE FROM logs")
	db.MustExec("DELETE FROM receipts")
	db.MustExec("DELETE FROM log_filters")
}

func NewTestDB(node core.Node) *DB {
	cfg, _ := config.NewConfig("private")
	db, _ := NewDB(cfg.Database, node)
	db.clearData()
	return db
}
