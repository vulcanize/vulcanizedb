package main

import (
	"flag"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "contract-hash=x1234")
	abiFilepath := flag.String("abi-filepath", "", "path/to/abifile.json")
	flag.Parse()

	contractAbiString := cmd.GetAbi(*abiFilepath, *contractHash)
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	watchedContract := core.Contract{
		Abi:  contractAbiString,
		Hash: *contractHash,
	}
	repository.CreateContract(watchedContract)
}
