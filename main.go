package main

import (
	"github.com/8thlight/vulcanizedb/core"
)

func main() {
	var blockchain core.Blockchain = core.NewGethBlockchain()
	blockchain.RegisterObserver(core.BlockchainLoggingObserver{})
	blockchain.SubscribeToEvents()
}
