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

	It("fills in the only missing block (Number 1)", func() {
		blocks := []core.Block{
			{Number: 1},
			{Number: 2},
		}
		blockchain := fakes.NewBlockchainWithBlocks(blocks)
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(core.Block{Number: 2})

		blocksAdded := history.PopulateMissingBlocks(blockchain, repository, 1)
		_, err := repository.FindBlockByNumber(1)

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
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(core.Block{Number: 1})
		repository.CreateOrUpdateBlock(core.Block{Number: 2})
		repository.CreateOrUpdateBlock(core.Block{Number: 3})
		repository.CreateOrUpdateBlock(core.Block{Number: 6})
		repository.CreateOrUpdateBlock(core.Block{Number: 7})
		repository.CreateOrUpdateBlock(core.Block{Number: 9})
		repository.CreateOrUpdateBlock(core.Block{Number: 11})
		repository.CreateOrUpdateBlock(core.Block{Number: 12})

		blocksAdded := history.PopulateMissingBlocks(blockchain, repository, 5)

		Expect(blocksAdded).To(Equal(3))
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

	It("returns the number of blocks created", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 6},
		})
		repository := repositories.NewInMemory()
		repository.CreateOrUpdateBlock(core.Block{Number: 3})
		repository.CreateOrUpdateBlock(core.Block{Number: 6})

		numberOfBlocksCreated := history.PopulateMissingBlocks(blockchain, repository, 3)

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
		repository := repositories.NewInMemory()

		history.RetrieveAndUpdateBlocks(blockchain, repository, history.MakeRange(2, 5))
		Expect(repository.BlockCount()).To(Equal(3))
		Expect(repository.CreateOrUpdateBlockCallCount).To(Equal(3))
	})

})
