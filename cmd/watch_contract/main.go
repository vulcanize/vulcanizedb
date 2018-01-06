package main

import (
	"flag"

	"github.com/vulcanize/vulcanizedb/cmd"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "contract-hash=x1234")
	abiFilepath := flag.String("abi-filepath", "", "path/to/abifile.json")
	flag.Parse()

	contractAbiString := cmd.GetAbi(*abiFilepath, *contractHash)
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	watchedContract := core.Contract{
		Abi:  contractAbiString,
		Hash: *contractHash,
	}
	repository.CreateContract(watchedContract)
}
