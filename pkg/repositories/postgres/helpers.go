package postgres

import (
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func ClearData(postgres *DB) {
	postgres.DB.MustExec("DELETE FROM watched_contracts")
	postgres.DB.MustExec("DELETE FROM transactions")
	postgres.DB.MustExec("DELETE FROM blocks")
	postgres.DB.MustExec("DELETE FROM logs")
	postgres.DB.MustExec("DELETE FROM receipts")
	postgres.DB.MustExec("DELETE FROM log_filters")
}

func BuildRepository(node core.Node) *DB {
	cfg, _ := config.NewConfig("private")
	repository, _ := NewDB(cfg.Database, node)
	ClearData(repository)
	return repository
}
