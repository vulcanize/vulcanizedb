package main

import (
	"log"

	"fmt"

	cfg "github.com/8thlight/vulcanizedb/config"
	"github.com/8thlight/vulcanizedb/core"
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
	port := config.Database.Port
	host := config.Database.Hostname
	databaseName := config.Database.Name

	var blockchain core.Blockchain = core.NewGethBlockchain(ipcPath)
	blockchain.RegisterObserver(core.BlockchainLoggingObserver{})
	pgConfig := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable", host, port, databaseName)
	db, err := sqlx.Connect("postgres", pgConfig)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	blockchain.RegisterObserver(core.BlockchainDBObserver{Db: db})
	blockchain.SubscribeToEvents()
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
