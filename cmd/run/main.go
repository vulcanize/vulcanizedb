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

	repository := repositories.NewPostgres(config.Database)
	fmt.Printf("Creating Geth Blockchain to: %s\n", config.Client.IPCPath)
	listener := blockchain_listener.NewBlockchainListener(
		geth.NewGethBlockchain(config.Client.IPCPath),
		[]core.BlockchainObserver{
			observers.BlockchainLoggingObserver{},
			observers.NewBlockchainDbObserver(repository),
		},
	)
	listener.Start()
}
