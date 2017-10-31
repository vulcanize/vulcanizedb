package integration_test

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/8thlight/vulcanizedb/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	_, filename, _, _ = runtime.Caller(0)
	basepath          = filepath.Dir(filename)
)

func RunTimePath() string {
	return path.Join(path.Dir(filename), "../")
}

type ObserverWithChannel struct {
	blocks chan core.Block
}

func (observer *ObserverWithChannel) NotifyBlockAdded(block core.Block) {
	fmt.Println("Block: ", block.Number)
	observer.blocks <- block
}

var _ = Describe("Reading from the Geth blockchain", func() {

	It("reads two blocks with incrementing numbers", func(done Done) {
		addedBlock := make(chan core.Block, 10)
		observer := &ObserverWithChannel{addedBlock}

		var blockchain core.Blockchain = core.NewGethBlockchain(RunTimePath() + "/test_data_dir/geth.ipc")
		blockchain.RegisterObserver(observer)

		go blockchain.SubscribeToEvents()

		firstBlock := <-addedBlock
		Expect(firstBlock).ShouldNot(BeNil())
		secondBlock := <-addedBlock
		Expect(firstBlock.Number + 1).Should(Equal(secondBlock.Number))

		close(done)
	}, 10)

})
