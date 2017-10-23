package core

import "fmt"

type BlockchainLoggingObserver struct{}

func (blockchainObserver BlockchainLoggingObserver) NotifyBlockAdded(block Block) {
	fmt.Println("Added block: %f", block.Number)
}
