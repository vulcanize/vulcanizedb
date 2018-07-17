package history

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

type BlockValidator struct {
	blockchain      core.Blockchain
	blockRepository datastore.BlockRepository
	windowSize      int
}

func NewBlockValidator(blockchain core.Blockchain, blockRepository datastore.BlockRepository, windowSize int) *BlockValidator {
	return &BlockValidator{
		blockchain:      blockchain,
		blockRepository: blockRepository,
		windowSize:      windowSize,
	}
}

func (bv BlockValidator) ValidateBlocks() ValidationWindow {
	window := MakeValidationWindow(bv.blockchain, bv.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	RetrieveAndUpdateBlocks(bv.blockchain, bv.blockRepository, blockNumbers)
	lastBlock := bv.blockchain.LastBlock().Int64()
	bv.blockRepository.SetBlocksStatus(lastBlock)
	return window
}
