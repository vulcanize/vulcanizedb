package main

import (
	"fmt"

	"github.com/8thlight/vulcanizedb/core"
)

func main() {
	fmt.Println("Starting connection")
	var blockchain core.Blockchain = core.NewGethBlockchain()
	blockchain.RegisterObserver(core.BlockchainLoggingObserver{})
}
