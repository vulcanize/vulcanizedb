package core

import "fmt"

type BlockchainLoggingObserver struct{}

func (blockchainObserver BlockchainLoggingObserver) NotifyBlockAdded(block Block) {
	fmt.Printf("New block was added: %d\n", block.Number)
}
