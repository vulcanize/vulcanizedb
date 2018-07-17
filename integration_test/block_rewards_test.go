package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Rewards calculations", func() {

	It("calculates a block reward for a real block", func() {
		blockchain := geth.NewBlockChain(test_config.InfuraClient.IPCPath)
		block, err := blockchain.GetBlockByNumber(1071819)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.Reward).To(Equal(5.31355))
	})

	It("calculates an uncle reward for a real block", func() {
		blockchain := geth.NewBlockChain(test_config.InfuraClient.IPCPath)
		block, err := blockchain.GetBlockByNumber(1071819)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.UnclesReward).To(Equal(6.875))
	})

})
