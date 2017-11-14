package main

import (
	"flag"

	"log"

	"fmt"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/8thlight/vulcanizedb/pkg/watched_contracts"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	flag.Parse()
	config := cmd.LoadConfig(*environment)

	repository := repositories.NewPostgres(config.Database)
	contractSummary, err := watched_contracts.NewSummary(repository, *contractHash)
	if err != nil {
		log.Fatalln(err)
	}
	output := watched_contracts.GenerateConsoleOutput(contractSummary)
	fmt.Println(output)
}
