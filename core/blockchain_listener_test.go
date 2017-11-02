package core_test

import (
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Blockchain listeners", func() {

	It("starts with no blocks", func(done Done) {
		observer := fakes.NewFakeBlockchainObserverTwo()
		blockchain := &fakes.Blockchain{}

		core.NewBlockchainListener(blockchain, []core.BlockchainObserver{observer})

		Expect(len(observer.CurrentBlocks)).To(Equal(0))
		close(done)
	}, 1)

	It("sees when one block was added", func(done Done) {
		observer := fakes.NewFakeBlockchainObserverTwo()
		blockchain := &fakes.Blockchain{}
		listener := core.NewBlockchainListener(blockchain, []core.BlockchainObserver{observer})
		go listener.Start()

		go blockchain.AddBlock(core.Block{Number: 123})

		wasObserverNotified := <-observer.WasNotified
		Expect(wasObserverNotified).To(BeTrue())
		Expect(len(observer.CurrentBlocks)).To(Equal(1))
		addedBlock := observer.CurrentBlocks[0]
		Expect(addedBlock.Number).To(Equal(int64(123)))
		close(done)
	}, 1)

	It("sees a second block", func(done Done) {
		observer := fakes.NewFakeBlockchainObserverTwo()
		blockchain := &fakes.Blockchain{}
		listener := core.NewBlockchainListener(blockchain, []core.BlockchainObserver{observer})
		go listener.Start()

		go blockchain.AddBlock(core.Block{Number: 123})
		<-observer.WasNotified
		go blockchain.AddBlock(core.Block{Number: 456})
		wasObserverNotified := <-observer.WasNotified

		Expect(wasObserverNotified).To(BeTrue())
		Expect(len(observer.CurrentBlocks)).To(Equal(2))
		addedBlock := observer.CurrentBlocks[1]
		Expect(addedBlock.Number).To(Equal(int64(456)))
		close(done)
	}, 1)

})
