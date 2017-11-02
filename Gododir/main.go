package main

import (
	"log"

	"fmt"

	"github.com/8thlight/vulcanizedb/blockchain_listener"
	cfg "github.com/8thlight/vulcanizedb/config"
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/geth"
	"github.com/8thlight/vulcanizedb/observers"
	"github.com/jmoiron/sqlx"
	do "gopkg.in/godo.v2"
)

func parseIpcPath(context *do.Context) string {
	ipcPath := context.Args.MayString("", "ipc-path", "i")
	if ipcPath == "" {
		log.Fatalln("--ipc-path required")
	}
	return ipcPath
}

func startBlockchainListener(config cfg.Config, ipcPath string) {
	blockchain := geth.NewGethBlockchain(ipcPath)
	loggingObserver := observers.BlockchainLoggingObserver{}
	connectString := cfg.DbConnectionString(cfg.Public().Database)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	dbObserver := (observers.BlockchainDBObserver{Db: db})
	listener := blockchain_listener.NewBlockchainListener(blockchain, []core.BlockchainObserver{
		loggingObserver,
		dbObserver,
	})
	listener.Start()
}

func tasks(p *do.Project) {

	p.Task("runPublic", nil, func(context *do.Context) {
		startBlockchainListener(cfg.Public(), parseIpcPath(context))
	})

	p.Task("runPrivate", nil, func(context *do.Context) {
		startBlockchainListener(cfg.Private(), parseIpcPath(context))
	})

	p.Task("migratePublic", nil, func(context *do.Context) {
		connectString := cfg.DbConnectionString(cfg.Public().Database)
		context.Bash(fmt.Sprintf("migrate -database '%s' -path ./migrations up", connectString))
		context.Bash(fmt.Sprintf("pg_dump -O -s %s > migrations/schema.sql", cfg.Public().Database.Name))
	})

	p.Task("migratePrivate", nil, func(context *do.Context) {
		connectString := cfg.DbConnectionString(cfg.Private().Database)
		context.Bash(fmt.Sprintf("migrate -database '%s' -path ./migrations up", connectString))
		context.Bash(fmt.Sprintf("pg_dump -O -s %s > migrations/schema.sql", cfg.Private().Database.Name))
	})

	p.Task("rollbackPublic", nil, func(context *do.Context) {
		connectString := cfg.DbConnectionString(cfg.Public().Database)
		context.Bash(fmt.Sprintf("migrate -database '%s' -path ./migrations down 1", connectString))
		context.Bash("pg_dump -O -s vulcanize_public > migrations/schema.sql")
	})

	p.Task("rollbackPrivate", nil, func(context *do.Context) {
		connectString := cfg.DbConnectionString(cfg.Private().Database)
		context.Bash(fmt.Sprintf("migrate -database '%s' -path ./migrations down 1", connectString))
		context.Bash("pg_dump -O -s vulcanize_private > migrations/schema.sql")
	})

}

func main() {
	do.Godo(tasks)
}
