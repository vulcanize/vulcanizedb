package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/filters"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	filterFilePath := flag.String("filter-filepath", "", "path/to/filter.json")

	flag.Parse()
	var logFilters filters.LogFilters
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	absFilePath := cmd.AbsFilePath(*filterFilePath)
	logFilterBytes, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(logFilterBytes, &logFilters)
	if err != nil {
		log.Fatal(err)
	}
	for _, filter := range logFilters {
		err = repository.AddFilter(filter)
		if err != nil {
			log.Fatal(err)
		}
	}
}
