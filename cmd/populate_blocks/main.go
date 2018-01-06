package main

import (
	"flag"

	"fmt"

	"github.com/vulcanize/vulcanizedb/cmd"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	startingBlockNumber := flag.Int("starting-number", -1, "First block to fill from")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	numberOfBlocksCreated := history.PopulateMissingBlocks(blockchain, repository, int64(*startingBlockNumber))
	fmt.Printf("Populated %d blocks", numberOfBlocksCreated)
}
