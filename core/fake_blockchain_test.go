package core_test

import (
	"math/big"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The fake blockchain", func() {

	It("conforms to the Blockchain interface", func() {
		var blockchain core.Blockchain = &fakes.Blockchain{}
		Expect(blockchain).ShouldNot(BeNil())
	})

	It("lets the only observer know when a block was added", func() {
		blockchain := fakes.Blockchain{}
		blockchainObserver := &fakes.BlockchainObserver{}
		blockchain.RegisterObserver(blockchainObserver)

		blockchain.AddBlock(core.Block{})

		Expect(blockchainObserver.WasToldBlockAdded()).Should(Equal(true))
	})

	It("lets the second observer know when a block was added", func() {
		blockchain := fakes.Blockchain{}
		blockchainObserverOne := &fakes.BlockchainObserver{}
		blockchainObserverTwo := &fakes.BlockchainObserver{}
		blockchain.RegisterObserver(blockchainObserverOne)
		blockchain.RegisterObserver(blockchainObserverTwo)

		blockchain.AddBlock(core.Block{})

		Expect(blockchainObserverTwo.WasToldBlockAdded()).Should(Equal(true))
	})

	It("passes the added block to the observer", func() {
		blockchain := fakes.Blockchain{}
		blockchainObserver := &fakes.BlockchainObserver{}
		blockchain.RegisterObserver(blockchainObserver)

		blockchain.AddBlock(core.Block{Number: big.NewInt(123)})

		Expect(blockchainObserver.LastAddedBlock().Number).ShouldNot(BeNil())
		Expect(blockchainObserver.LastAddedBlock().Number).Should(Equal(big.NewInt(123)))
	})

})
