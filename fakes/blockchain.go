package fakes

import "github.com/8thlight/vulcanizedb/core"

type Blockchain struct {
	outputBlocks  chan core.Block
	WasToldToStop bool
}

func (blockchain *Blockchain) SubscribeToBlocks(outputBlocks chan core.Block) {
	blockchain.outputBlocks = outputBlocks
}

func (blockchain Blockchain) AddBlock(block core.Block) {
	blockchain.outputBlocks <- block
}

func (*Blockchain) StartListening() {}

func (blockchain *Blockchain) StopListening() {
	blockchain.WasToldToStop = true
}
