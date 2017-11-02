package fakes

import "github.com/8thlight/vulcanizedb/core"

type Blockchain struct {
	outputBlocks chan core.Block
}

func (blockchain *Blockchain) SubscribeToBlocks(outputBlocks chan core.Block) {
	blockchain.outputBlocks = outputBlocks
}

func (blockchain Blockchain) AddBlock(block core.Block) {
	blockchain.outputBlocks <- block
}

func (*Blockchain) StartListening() {}
