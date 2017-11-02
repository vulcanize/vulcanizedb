package main

import (
	"log"

	"fmt"

	"github.com/8thlight/vulcanizedb/blockchain_listener"
	"github.com/8thlight/vulcanizedb/config"
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/geth"
	"github.com/8thlight/vulcanizedb/observers"
	"github.com/jmoiron/sqlx"
	do "gopkg.in/godo.v2"
)

func parseEnvironment(context *do.Context) string {
	environment := context.Args.MayString("", "environment", "env", "e")
	if environment == "" {
		log.Fatalln("--environment required")
	}
	return environment
}

func startBlockchainListener(cfg config.Config) {
	fmt.Println("Client Path ", cfg.Client.IPCPath)
	blockchain := geth.NewGethBlockchain(cfg.Client.IPCPath)
	loggingObserver := observers.BlockchainLoggingObserver{}
	connectString := config.DbConnectionString(cfg.Database)
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

	p.Task("run", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := config.NewConfig(environment)
		startBlockchainListener(cfg)
	})

	p.Task("migrate", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := config.NewConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./migrations up", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > migrations/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

	p.Task("rollback", nil, func(context *do.Context) {
		environment := parseEnvironment(context)
		cfg := config.NewConfig(environment)
		connectString := config.DbConnectionString(cfg.Database)
		migrate := fmt.Sprintf("migrate -database '%s' -path ./migrations down 1", connectString)
		dumpSchema := fmt.Sprintf("pg_dump -O -s %s > migrations/schema.sql", cfg.Database.Name)
		context.Bash(migrate)
		context.Bash(dumpSchema)
	})

}

func main() {
	do.Godo(tasks)
}
