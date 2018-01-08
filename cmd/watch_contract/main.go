package main

import (
	"flag"

	"strings"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "contract-hash=x1234")
	abiFilepath := flag.String("abi-filepath", "", "path/to/abifile.json")
	network := flag.String("network", "", "ropsten")

	flag.Parse()
	contractHashLowered := strings.ToLower(*contractHash)

	contractAbiString := cmd.GetAbi(*abiFilepath, contractHashLowered, *network)
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	watchedContract := core.Contract{
		Abi:  contractAbiString,
		Hash: contractHashLowered,
	}
	repository.CreateContract(watchedContract)
}
