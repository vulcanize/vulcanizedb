package blockchain_listener

import "github.com/8thlight/vulcanizedb/core"

type BlockchainListener struct {
	inputBlocks chan core.Block
	blockchain  core.Blockchain
	observers   []core.BlockchainObserver
}

func NewBlockchainListener(blockchain core.Blockchain, observers []core.BlockchainObserver) BlockchainListener {
	inputBlocks := make(chan core.Block, 10)
	blockchain.SubscribeToBlocks(inputBlocks)
	listener := BlockchainListener{
		inputBlocks: inputBlocks,
		blockchain:  blockchain,
		observers:   observers,
	}
	return listener
}

func (listener BlockchainListener) Start() {
	go listener.blockchain.StartListening()
	for block := range listener.inputBlocks {
		listener.notifyObservers(block)
	}
}

func (listener BlockchainListener) notifyObservers(block core.Block) {
	for _, observer := range listener.observers {
		observer.NotifyBlockAdded(block)
	}
}

func (listener BlockchainListener) Stop() {
	listener.blockchain.StopListening()
}
