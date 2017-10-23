package main

import (
	"flag"
	"github.com/8thlight/vulcanizedb/core"
)

func main() {
	ipcPath := flag.String("ipcPath", "", "location geth.ipc")
	flag.Parse()

	var blockchain core.Blockchain = core.NewGethBlockchain(*ipcPath)
	blockchain.RegisterObserver(core.BlockchainLoggingObserver{})
	blockchain.SubscribeToEvents()
}
