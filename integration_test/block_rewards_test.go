package integration

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cfg "github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

var _ = Describe("Rewards calculations", func() {

	It("calculates a block reward for a real block", func() {
		config, err := cfg.NewConfig("infura")
		if err != nil {
			log.Fatalln(err)
		}
		blockchain := geth.NewBlockchain(config.Client.IPCPath)
		block := blockchain.GetBlockByNumber(1071819)
		Expect(block.Reward).To(Equal(5.31355))
	})

	It("calculates an uncle reward for a real block", func() {
		config, err := cfg.NewConfig("infura")
		if err != nil {
			log.Fatalln(err)
		}
		blockchain := geth.NewBlockchain(config.Client.IPCPath)
		block := blockchain.GetBlockByNumber(1071819)
		Expect(block.UnclesReward).To(Equal(6.875))
	})

})
