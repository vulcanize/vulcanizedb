package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"math/big"
)

var _ = Describe("Blocks validator", func() {

	It("calls create or update for all blocks within the window", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(7))
		blocksRepository := fakes.NewMockBlockRepository()
		validator := history.NewBlockValidator(blockChain, blocksRepository, 2)

		window := validator.ValidateBlocks()

		Expect(window).To(Equal(history.ValidationWindow{LowerBound: 5, UpperBound: 7}))
		blocksRepository.AssertCreateOrUpdateBlockCallCountEquals(3)
	})

	It("returns the number of largest block", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(3))
		maxBlockNumber := blockChain.LastBlock()

		Expect(maxBlockNumber.Int64()).To(Equal(int64(3)))
	})
})
