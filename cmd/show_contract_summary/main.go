package main

import (
	"flag"

	"log"

	"fmt"

	"math/big"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/contract_summary"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	_blockNumber := flag.Int64("block-number", -1, "Block number of summary")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	repository := cmd.LoadPostgres(config.Database, blockchain.Node())
	blockNumber := requestedBlockNumber(_blockNumber)

	contractSummary, err := contract_summary.NewSummary(blockchain, repository, *contractHash, blockNumber)
	if err != nil {
		log.Fatalln(err)
	}
	output := contract_summary.GenerateConsoleOutput(contractSummary)
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
