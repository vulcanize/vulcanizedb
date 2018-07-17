package history_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Populating blocks", func() {
	var inMemory *inmemory.InMemory
	var blockRepository *inmemory.BlockRepository

	BeforeEach(func() {
		inMemory = inmemory.NewInMemory()
		blockRepository = &inmemory.BlockRepository{InMemory: inMemory}
	})

	It("fills in the only missing block (BlockNumber 1)", func() {
		blocks := []core.Block{
			{Number: 1},
			{Number: 2},
		}
		blockchain := fakes.NewBlockchainWithBlocks(blocks)

		blockRepository.CreateOrUpdateBlock(core.Block{Number: 2})

		blocksAdded := history.PopulateMissingBlocks(blockchain, blockRepository, 1)
		_, err := blockRepository.GetBlock(1)

		Expect(blocksAdded).To(Equal(1))
		Expect(err).ToNot(HaveOccurred())
	})

	It("fills in the three missing blocks (Numbers: 5,8,10)", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
			{Number: 7},
			{Number: 8},
			{Number: 9},
			{Number: 10},
			{Number: 11},
			{Number: 12},
			{Number: 13},
		})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 1})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 2})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 3})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 6})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 7})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 9})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 11})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 12})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 13})

		blocksAdded := history.PopulateMissingBlocks(blockchain, blockRepository, 5)

		Expect(blocksAdded).To(Equal(3))
		Expect(blockRepository.BlockCount()).To(Equal(12))
		_, err := blockRepository.GetBlock(4)
		Expect(err).To(HaveOccurred())
		_, err = blockRepository.GetBlock(5)
		Expect(err).ToNot(HaveOccurred())
		_, err = blockRepository.GetBlock(8)
		Expect(err).ToNot(HaveOccurred())
		_, err = blockRepository.GetBlock(10)
		Expect(err).ToNot(HaveOccurred())
		_, err = blockRepository.GetBlock(14)
		Expect(err).To(HaveOccurred())
	})

	It("returns the number of blocks created", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
		})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 3})
		blockRepository.CreateOrUpdateBlock(core.Block{Number: 6})

		numberOfBlocksCreated := history.PopulateMissingBlocks(blockchain, blockRepository, 3)

		Expect(numberOfBlocksCreated).To(Equal(2))
	})

	It("updates the repository with a range of blocks w/in the range ", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
			{Number: 4},
			{Number: 5},
		})

		history.RetrieveAndUpdateBlocks(blockchain, blockRepository, history.MakeRange(2, 5))
		Expect(blockRepository.BlockCount()).To(Equal(4))
		Expect(blockRepository.CreateOrUpdateBlockCallCount).To(Equal(4))
	})

	It("does not call repository create block when there is an error", func() {
		blockchain := fakes.NewBlockchain(errors.New("error getting block"))
		blocks := history.MakeRange(1, 10)
		history.RetrieveAndUpdateBlocks(blockchain, blockRepository, blocks)
		Expect(blockRepository.BlockCount()).To(Equal(0))
		Expect(blockRepository.CreateOrUpdateBlockCallCount).To(Equal(0))
	})
})
