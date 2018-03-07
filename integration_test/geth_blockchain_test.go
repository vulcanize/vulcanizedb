package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Reading from the Geth blockchain", func() {

	var blockchain *geth.Blockchain
	var inMemory *inmemory.InMemory

	BeforeEach(func() {
		blockchain = geth.NewBlockchain(test_config.TestClientConfig.IPCPath)
		inMemory = inmemory.NewInMemory()
	})

	It("reads two blocks", func(done Done) {
		blocks := &inmemory.BlockRepository{InMemory: inMemory}
		validator := history.NewBlockValidator(blockchain, blocks, 2)
		validator.ValidateBlocks()
		Expect(blocks.BlockCount()).To(Equal(2))
		close(done)
	}, 15)

	It("retrieves the genesis block and first block", func(done Done) {
		genesisBlock := blockchain.GetBlockByNumber(int64(0))
		firstBlock := blockchain.GetBlockByNumber(int64(1))
		lastBlockNumber := blockchain.LastBlock()

		Expect(genesisBlock.Number).To(Equal(int64(0)))
		Expect(firstBlock.Number).To(Equal(int64(1)))
		Expect(lastBlockNumber.Int64()).To(BeNumerically(">", 0))
		close(done)
	}, 15)

	It("retrieves the node info", func(done Done) {
		node := blockchain.Node()
		devNetworkNodeId := float64(1)

		Expect(node.GenesisBlock).ToNot(BeNil())
		Expect(node.NetworkID).To(Equal(devNetworkNodeId))
		Expect(len(node.ID)).To(Equal(128))
		Expect(node.ClientName).To(ContainSubstring("Geth"))

		close(done)
	}, 15)

})
