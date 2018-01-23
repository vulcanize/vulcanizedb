package main

import (
	"flag"

	"time"

	"os"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

const (
	pollingInterval = 7 * time.Second
)

func backFillAllBlocks(blockchain core.Blockchain, repository repositories.Postgres, missingBlocksPopulated chan int, startingBlockNumber int64) {
	go func() {
		missingBlocksPopulated <- history.PopulateMissingBlocks(blockchain, repository, startingBlockNumber)
	}()
}

func main() {
	environment := flag.String("environment", "", "Environment name")
	startingBlockNumber := flag.Int("starting-number", 0, "First block to fill from")
	flag.Parse()

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	validator := history.NewBlockValidator(blockchain, repository, 15)

	missingBlocksPopulated := make(chan int)
	_startingBlockNumber := int64(*startingBlockNumber)
	go backFillAllBlocks(blockchain, repository, missingBlocksPopulated, _startingBlockNumber)

	for {
		select {
		case <-ticker.C:
			window := validator.ValidateBlocks()
			validator.Log(os.Stdout, window)
		case <-missingBlocksPopulated:
			go backFillAllBlocks(blockchain, repository, missingBlocksPopulated, _startingBlockNumber)
		}
	}
}
