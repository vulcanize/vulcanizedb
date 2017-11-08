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
		repository.CreateBlock(core.Block{Number: 2})

		history.PopulateBlocks(blockchain, repository, 1)

		block := repository.FindBlockByNumber(1)
		Expect(block).NotTo(BeNil())
		Expect(block.Hash).To(Equal("x012343"))
	})

	It("fills in two missing blocks", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
			{Number: 8},
			{Number: 10},
			{Number: 13},
		})
		repository := repositories.NewInMemory()
		repository.CreateBlock(core.Block{Number: 1})
		repository.CreateBlock(core.Block{Number: 2})
		repository.CreateBlock(core.Block{Number: 3})
		repository.CreateBlock(core.Block{Number: 6})
		repository.CreateBlock(core.Block{Number: 7})
		repository.CreateBlock(core.Block{Number: 9})
		repository.CreateBlock(core.Block{Number: 11})
		repository.CreateBlock(core.Block{Number: 12})

		history.PopulateBlocks(blockchain, repository, 5)

		Expect(repository.BlockCount()).To(Equal(11))
		Expect(repository.FindBlockByNumber(4)).To(BeNil())
		Expect(repository.FindBlockByNumber(5)).NotTo(BeNil())
		Expect(repository.FindBlockByNumber(8)).NotTo(BeNil())
		Expect(repository.FindBlockByNumber(10)).NotTo(BeNil())
		Expect(repository.FindBlockByNumber(13)).To(BeNil())
	})

	It("returns the number of blocks created", func() {
		blockchain := fakes.NewBlockchainWithBlocks([]core.Block{
			{Number: 4},
			{Number: 5},
		})
		repository := repositories.NewInMemory()
		repository.CreateBlock(core.Block{Number: 3})
		repository.CreateBlock(core.Block{Number: 6})

		numberOfBlocksCreated := history.PopulateBlocks(blockchain, repository, 3)

		Expect(numberOfBlocksCreated).To(Equal(2))
	})

})
