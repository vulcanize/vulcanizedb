package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Blocks validator", func() {

	It("calls create or update for all blocks within the window", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
			{Number: 7},
		})
		inMemoryDB := inmemory.NewInMemory()
		blocksRepository := &inmemory.BlockRepository{InMemory: inMemoryDB}
		validator := history.NewBlockValidator(blockchain, blocksRepository, 2)

		window := validator.ValidateBlocks()

		Expect(window).To(Equal(history.ValidationWindow{LowerBound: 5, UpperBound: 7}))
		Expect(blocksRepository.BlockCount()).To(Equal(3))
		Expect(blocksRepository.CreateOrUpdateBlockCallCount).To(Equal(3))
	})

	It("returns the number of largest block", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
		})
		maxBlockNumber := blockchain.LastBlock()

		Expect(maxBlockNumber.Int64()).To(Equal(int64(3)))
	})
})
