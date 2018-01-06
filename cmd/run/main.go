package main

import (
	"flag"

	"time"

	"os"

	"github.com/vulcanize/vulcanizedb/cmd"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

const (
	pollingInterval = 7 * time.Second
)

func main() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	environment := flag.String("environment", "", "Environment name")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	validator := history.NewBlockValidator(blockchain, repository, 15)

	for range ticker.C {
		window := validator.ValidateBlocks()
		validator.Log(os.Stdout, window)
	}
}
