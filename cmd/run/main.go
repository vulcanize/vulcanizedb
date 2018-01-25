package main

import (
	"flag"

	"time"

	"os"

	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/utils"
)

const (
	pollingInterval = 7 * time.Second
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	flag.Parse()

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	config := utils.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := utils.LoadPostgres(config.Database, blockchain.Node())
	validator := history.NewBlockValidator(blockchain, repository, 15)

	for range ticker.C {
		window := validator.ValidateBlocks()
		validator.Log(os.Stdout, window)
	}
}
