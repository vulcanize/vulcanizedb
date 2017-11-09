package integration_test

import (
	"github.com/8thlight/vulcanizedb/pkg/blockchain_listener"
	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/fakes"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reading from the Geth blockchain", func() {

	var listener blockchain_listener.BlockchainListener
	var observer *fakes.BlockchainObserver
	var blockchain *geth.GethBlockchain

	BeforeEach(func() {
		observer = fakes.NewFakeBlockchainObserver()
		cfg, _ := config.NewConfig("private")
		blockchain = geth.NewGethBlockchain(cfg.Client.IPCPath)
		observers := []core.BlockchainObserver{observer}
		listener = blockchain_listener.NewBlockchainListener(blockchain, observers)
	})

	AfterEach(func() {
		listener.Stop()
	})

	It("reads two blocks", func(done Done) {
		go listener.Start()

		<-observer.WasNotified
		firstBlock := observer.LastBlock()
		Expect(firstBlock).NotTo(BeNil())

		<-observer.WasNotified
		secondBlock := observer.LastBlock()
		Expect(secondBlock).NotTo(BeNil())

		Expect(firstBlock.Number + 1).Should(Equal(secondBlock.Number))

		close(done)
	}, 10)

	It("retrieves the genesis block and first block", func(done Done) {
		genesisBlock := blockchain.GetBlockByNumber(int64(0))
		firstBlock := blockchain.GetBlockByNumber(int64(1))

		Expect(genesisBlock.Number).To(Equal(int64(0)))
		Expect(firstBlock.Number).To(Equal(int64(1)))

		close(done)
	}, 10)

})
