package main

import (
	"flag"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "contract-hash=x1234")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	repository := cmd.LoadPostgres(config.Database)
	repository.CreateWatchedContract(core.WatchedContract{Hash: *contractHash})
}
