package repositories

import (
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func ClearData(postgres Postgres) {
	postgres.Db.MustExec("DELETE FROM watched_contracts")
	postgres.Db.MustExec("DELETE FROM transactions")
	postgres.Db.MustExec("DELETE FROM blocks")
	postgres.Db.MustExec("DELETE FROM logs")
	postgres.Db.MustExec("DELETE FROM receipts")
	postgres.Db.MustExec("DELETE FROM log_filters")
}

func BuildRepository(node core.Node) Repository {
	cfg, _ := config.NewConfig("private")
	repository, _ := NewPostgres(cfg.Database, node)
	ClearData(repository)
	return repository
}
