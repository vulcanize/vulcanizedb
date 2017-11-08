package main

import (
	"flag"

	"log"

	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/jmoiron/sqlx"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	startingBlockNumber := flag.Int("starting-number", -1, "First block to fill from")
	flag.Parse()
	cfg := config.NewConfig(*environment)

	blockchain := geth.NewGethBlockchain(cfg.Client.IPCPath)
	connectString := config.DbConnectionString(cfg.Database)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	repository := repositories.NewPostgres(db)
	numberOfBlocksCreated := history.PopulateBlocks(blockchain, repository, int64(*startingBlockNumber))
	fmt.Printf("Populated %d blocks", numberOfBlocksCreated)
}
