package main

import (
	"fmt"
	"log"

	"flag"

	"github.com/8thlight/vulcanizedb/pkg/blockchain_listener"
	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/observers"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/jmoiron/sqlx"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	flag.Parse()
	cfg := config.NewConfig(*environment)

	fmt.Println("Client Path ", cfg.Client.IPCPath)
	blockchain := geth.NewGethBlockchain(cfg.Client.IPCPath)
	loggingObserver := observers.BlockchainLoggingObserver{}
	connectString := config.DbConnectionString(cfg.Database)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	repository := repositories.NewPostgres(db)
	dbObserver := observers.NewBlockchainDbObserver(repository)
	listener := blockchain_listener.NewBlockchainListener(blockchain, []core.BlockchainObserver{
		loggingObserver,
		dbObserver,
	})
	listener.Start()
}
