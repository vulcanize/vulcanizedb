package integration

import (
	"log"

	cfg "github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reading contracts", func() {

	Describe("Block and Uncle rewards", func() {
		It("calculates a block reward for a real block", func() {
			config, err := cfg.NewConfig("infura")
			if err != nil {
				log.Fatalln(err)
			}
			blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
			block := blockchain.GetBlockByNumber(1071819)
			Expect(block.BlockReward).To(Equal(5.31355))
		})

		It("calculates an uncle reward for a real block", func() {
			config, err := cfg.NewConfig("infura")
			if err != nil {
				log.Fatalln(err)
			}
			blockchain := geth.NewGethBlockchain(config.Client.IPCPath)
			block := blockchain.GetBlockByNumber(1071819)
			Expect(block.UncleReward).To(Equal(6.875))
		})
	})

})
