package history_test

import (
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/fakes"
	"github.com/8thlight/vulcanizedb/pkg/history"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Populating blocks", func() {

	It("fills in the only missing block", func() {
		blocks := []core.Block{{Number: 1, Hash: "x012343"}}
		blockchain := fakes.NewBlockchainWithBlocks(blocks)
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(core.Block{Number: 2})

		history.PopulateMissingBlocks(blockchain, repository, 1)

		block, err := repository.FindBlockByNumber(1)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.Hash).To(Equal("x012343"))
	})

	It("fills in the three missing blocks (5,8,10)", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 8},
			{Number: 10},
			{Number: 13},
		})
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(core.Block{Number: 1})
		repository.CreateOrUpdateBlock(core.Block{Number: 2})
		repository.CreateOrUpdateBlock(core.Block{Number: 3})
		repository.CreateOrUpdateBlock(core.Block{Number: 6})
		repository.CreateOrUpdateBlock(core.Block{Number: 7})
		repository.CreateOrUpdateBlock(core.Block{Number: 9})
		repository.CreateOrUpdateBlock(core.Block{Number: 11})
		repository.CreateOrUpdateBlock(core.Block{Number: 12})

		history.PopulateMissingBlocks(blockchain, repository, 5)

		Expect(repository.BlockCount()).To(Equal(11))
		_, err := repository.FindBlockByNumber(4)
		Expect(err).To(HaveOccurred())
		_, err = repository.FindBlockByNumber(5)
		Expect(err).ToNot(HaveOccurred())
		_, err = repository.FindBlockByNumber(8)
		Expect(err).ToNot(HaveOccurred())
		_, err = repository.FindBlockByNumber(10)
		Expect(err).ToNot(HaveOccurred())
		_, err = repository.FindBlockByNumber(13)
		Expect(err).To(HaveOccurred())
	})

	It("updates the repository with a range of blocks w/in sliding window ", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 1},
			{Number: 2},
			{Number: 3},
			{Number: 4},
			{Number: 5},
		})
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(blockchain.GetBlockByNumber(5))

		history.UpdateBlocksWindow(blockchain, repository, 2)

		Expect(repository.BlockCount()).To(Equal(3))
		Expect(repository.HandleBlockCallCount).To(Equal(3))
	})

	It("Generates a range of int64", func() {
		numberOfBlocksCreated := history.MakeRange(0, 5)
		expected := []int64{0, 1, 2, 3, 4}

		Expect(numberOfBlocksCreated).To(Equal(expected))
	})

	It("returns the number of blocks created", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
		})
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(core.Block{Number: 3})
		repository.CreateOrUpdateBlock(core.Block{Number: 6})

		numberOfBlocksCreated := history.PopulateMissingBlocks(blockchain, repository, 3)

		Expect(numberOfBlocksCreated).To(Equal(2))
	})

	It("returns the window size", func() {
		window := history.Window{1, 3, 10}
		Expect(window.Size()).To(Equal(2))
	})

})
