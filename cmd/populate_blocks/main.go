package main

import (
	"flag"

	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/utils"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	startingBlockNumber := flag.Int("starting-number", -1, "First block to fill from")
	flag.Parse()
	config := utils.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := utils.LoadPostgres(config.Database, blockchain.Node())
	numberOfBlocksCreated := history.PopulateMissingBlocks(blockchain, repository, int64(*startingBlockNumber))
	fmt.Printf("Populated %d blocks", numberOfBlocksCreated)
}
