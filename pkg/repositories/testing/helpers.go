package testing

import "github.com/8thlight/vulcanizedb/pkg/repositories"

func ClearData(postgres repositories.Postgres) {
	postgres.Db.MustExec("DELETE FROM watched_contracts")
	postgres.Db.MustExec("DELETE FROM transactions")
	postgres.Db.MustExec("DELETE FROM blocks")
}
