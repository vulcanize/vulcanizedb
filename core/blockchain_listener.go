package core

type BlockchainListener struct {
	inputBlocks chan Block
	blockchain  Blockchain
	observers   []BlockchainObserver
}

func NewBlockchainListener(blockchain Blockchain, observers []BlockchainObserver) BlockchainListener {
	inputBlocks := make(chan Block, 10)
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

func (listener BlockchainListener) notifyObservers(block Block) {
	for _, observer := range listener.observers {
		observer.NotifyBlockAdded(block)
	}
}
