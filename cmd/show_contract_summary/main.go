package main

import (
	"flag"

	"log"

	"fmt"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/watched_contracts"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database)
	contractSummary, err := watched_contracts.NewSummary(blockchain, repository, *contractHash)
	if err != nil {
		log.Fatalln(err)
	}
	output := watched_contracts.GenerateConsoleOutput(contractSummary)
	fmt.Println(output)
}
