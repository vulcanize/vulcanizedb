package main

import (
	"flag"

	"log"

	"fmt"

	"math/big"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/8thlight/vulcanizedb/pkg/watched_contracts"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	_blockNumber := flag.Int64("block-number", -1, "Block number of summary")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database)
	blockNumber := requestedBlockNumber(_blockNumber)

	contractSummary, err := watched_contracts.NewSummary(blockchain, repository, *contractHash, blockNumber)
	if err != nil {
		log.Fatalln(err)
	}
	output := watched_contracts.GenerateConsoleOutput(contractSummary)
	fmt.Println(output)
}

func requestedBlockNumber(blockNumber *int64) *big.Int {
	var _blockNumber *big.Int
	if *blockNumber == -1 {
		_blockNumber = nil
	} else {
		_blockNumber = big.NewInt(*blockNumber)
	}
	return _blockNumber
}
