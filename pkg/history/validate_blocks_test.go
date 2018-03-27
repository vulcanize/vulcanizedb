package history_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Blocks validator", func() {

	It("creates a ValidationWindow equal to (HEAD-windowSize, HEAD)", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
			{Number: 4},
			{Number: 5},
		})

		validationWindow := history.MakeValidationWindow(blockchain, 2)

		Expect(validationWindow.LowerBound).To(Equal(int64(3)))
		Expect(validationWindow.UpperBound).To(Equal(int64(5)))
	})

	It("returns the window size", func() {
		window := history.ValidationWindow{LowerBound: 1, UpperBound: 3}
		Expect(window.Size()).To(Equal(2))
	})

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
		Expect(blocksRepository.BlockCount()).To(Equal(2))
		Expect(blocksRepository.CreateOrUpdateBlockCallCount).To(Equal(2))
	})

	It("logs window message", func() {
		inMemoryDB := inmemory.NewInMemory()
		blockRepository := &inmemory.BlockRepository{InMemory: inMemoryDB}

		expectedMessage := &bytes.Buffer{}
		window := history.ValidationWindow{LowerBound: 5, UpperBound: 7}
		history.ParsedWindowTemplate.Execute(expectedMessage, history.ValidationWindow{LowerBound: 5, UpperBound: 7})

		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{})
		validator := history.NewBlockValidator(blockchain, blockRepository, 2)
		actualMessage := &bytes.Buffer{}
		validator.Log(actualMessage, window)
		Expect(actualMessage).To(Equal(expectedMessage))
	})

	It("generates a range of int64s", func() {
		numberOfBlocksCreated := history.MakeRange(0, 5)
		expected := []int64{0, 1, 2, 3, 4}

		Expect(numberOfBlocksCreated).To(Equal(expected))
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
