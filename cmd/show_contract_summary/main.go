package main

import (
	"flag"

	"log"

	"fmt"

	"strings"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/contract_summary"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	_blockNumber := flag.Int64("block-number", -1, "Block number of summary")
	flag.Parse()

	contractHashLowered := strings.ToLower(*contractHash)
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	blockNumber := cmd.RequestedBlockNumber(_blockNumber)

	contractSummary, err := contract_summary.NewSummary(blockchain, repository, contractHashLowered, blockNumber)
	if err != nil {
		log.Fatalln(err)
	}
	output := contract_summary.GenerateConsoleOutput(contractSummary)
	fmt.Println(output)
}
