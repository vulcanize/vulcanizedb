package main

import (
	"fmt"
	"log"

	"flag"

	"github.com/8thlight/vulcanizedb/cmd"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth"
)

func main() {
	environment := flag.String("environment", "", "Environment name")
	contractHash := flag.String("contract-hash", "", "Contract hash to show summary")
	_blockNumber := flag.Int64("block-number", -1, "Block number of summary")
	flag.Parse()
	config := cmd.LoadConfig(*environment)
	blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
	blockNumber := cmd.RequestedBlockNumber(_blockNumber)

	logs, err := blockchain.GetLogs(core.Contract{Hash: *contractHash}, blockNumber)
	if err != nil {
		log.Fatalln(err)
	}
	for _, l := range logs {
		fmt.Println("\tAddress: ", l.Address)
		fmt.Println("\tTxHash: ", l.TxHash)
		fmt.Println("\tBlockNumber ", l.BlockNumber)
		fmt.Println("\tTopics: ")
		for i, topic := range l.Topics {
			fmt.Printf("\t\tTopic %d: %s\n", i, topic)
		}
		fmt.Printf("\tData: %s", l.Data)
		fmt.Print("\n\n")

	}

}
