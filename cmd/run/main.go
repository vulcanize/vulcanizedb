package main

import (
	"fmt"

	"flag"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/blockchain_listener"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/observers"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	flag.Parse()
	config := cmd.LoadConfig(*environment)

	fmt.Println("Client Path ", config.Client.IPCPath)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	loggingObserver := observers.BlockchainLoggingObserver{}
	repository := repositories.NewPostgres(config.Database)
	dbObserver := observers.NewBlockchainDbObserver(repository)
	listener := blockchain_listener.NewBlockchainListener(blockchain, []core.BlockchainObserver{
		loggingObserver,
		dbObserver,
	})
	listener.Start()
}
