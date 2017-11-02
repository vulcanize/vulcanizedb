package integration_test

import (
	"path"
	"runtime"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	_, filename, _, _ = runtime.Caller(0)
)

func RunTimePath() string {
	return path.Join(path.Dir(filename), "../")
}

var _ = Describe("Reading from the Geth blockchain", func() {

	It("reads two block with listener", func(done Done) {
		observer := fakes.NewFakeBlockchainObserverTwo()
		blockchain := core.NewGethBlockchain(RunTimePath() + "/test_data_dir/geth.ipc")
		observers := []core.BlockchainObserver{observer}
		listener := core.NewBlockchainListener(blockchain, observers)
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

})
