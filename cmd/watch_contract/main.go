package main

import (
	"flag"

	"strings"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/utils"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "contract-hash=x1234")
	abiFilepath := flag.String("abi-filepath", "", "path/to/abifile.json")
	network := flag.String("network", "", "ropsten")

	flag.Parse()
	contractHashLowered := strings.ToLower(*contractHash)

	contractAbiString := utils.GetAbi(*abiFilepath, contractHashLowered, *network)
	config := utils.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := utils.LoadPostgres(config.Database, blockchain.Node())
	watchedContract := core.Contract{
		Abi:  contractAbiString,
		Hash: contractHashLowered,
	}
	repository.CreateContract(watchedContract)
}
