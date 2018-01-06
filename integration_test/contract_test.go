package integration

import (
	"math/big"

	"log"

	cfg "github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reading contracts", func() {

	//TODO was experiencing Infura issue (I suspect) on 1/5. Unignore these and revisit if persists on next commit
	XDescribe("Reading the list of attributes", func() {
		It("returns a string attribute for a real contract", func() {
			config, err := cfg.NewConfig("infura")
			if err != nil {
				log.Fatalln(err)
			}
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			contractAttributes, err := blockchain.GetAttributes(contract)

			Expect(err).To(BeNil())
			Expect(len(contractAttributes)).NotTo(Equal(0))
			symbolAttribute := *testing.FindAttribute(contractAttributes, "symbol")
			Expect(symbolAttribute.Name).To(Equal("symbol"))
			Expect(symbolAttribute.Type).To(Equal("string"))
		})

		It("does not return an attribute that takes an input", func() {
			config, err := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			contractAttributes, err := blockchain.GetAttributes(contract)

			Expect(err).To(BeNil())
			attribute := testing.FindAttribute(contractAttributes, "balanceOf")
			Expect(attribute).To(BeNil())
		})

		It("does not return an attribute that is not constant", func() {
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			contractAttributes, err := blockchain.GetAttributes(contract)

			Expect(err).To(BeNil())
			attribute := testing.FindAttribute(contractAttributes, "unpause")
			Expect(attribute).To(BeNil())
		})
	})

	//TODO was experiencing Infura issue (I suspect) on 1/5. Unignore these and revisit if persists on next commit
	XDescribe("Getting a contract attribute", func() {
		It("returns the correct attribute for a real contract", func() {
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)

			contract := testing.SampleContract()
			name, err := blockchain.GetAttribute(contract, "name", nil)

			Expect(err).To(BeNil())
			Expect(name).To(Equal("OMGToken"))
		})

		It("returns the correct attribute for a real contract", func() {
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			name, err := blockchain.GetAttribute(contract, "name", nil)

			Expect(err).To(BeNil())
			Expect(name).To(Equal("OMGToken"))
		})

		It("returns the correct attribute for a real contract at a specific block height", func() {
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			name, err := blockchain.GetAttribute(contract, "name", big.NewInt(4701536))

			Expect(name).To(Equal("OMGToken"))
			Expect(err).To(BeNil())
		})

		It("returns an error when asking for an attribute that does not exist", func() {
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			name, err := blockchain.GetAttribute(contract, "missing_attribute", nil)

			Expect(err).To(Equal(geth.ErrInvalidStateAttribute))
			Expect(name).To(BeNil())
		})

		It("retrieves the event log for a specific block and contract", func() {
			expectedLogZero := core.Log{
				BlockNumber: 4703824,
				TxHash:      "0xf896bfd1eb539d881a1a31102b78de9f25cd591bf1fe1924b86148c0b205fd5d",
				Address:     "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
				Topics: map[int]string{
					0: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					1: "0x000000000000000000000000fbb1b73c4f0bda4f67dca266ce6ef42f520fbb98",
					2: "0x000000000000000000000000d26114cd6ee289accf82350c8d8487fedb8a0c07",
				},
				Index: 19,
				Data:  "0x0000000000000000000000000000000000000000000000000c7d713b49da0000"}
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)
			contract := testing.SampleContract()

			logs, err := blockchain.GetLogs(contract, big.NewInt(4703824), nil)

			Expect(err).To(BeNil())
			Expect(len(logs)).To(Equal(3))
			Expect(logs[0]).To(Equal(expectedLogZero))

		})

		It("returns and empty log array when no events for a given block / contract combo", func() {
			config, _ := cfg.NewConfig("infura")
			blockchain := geth.NewBlockchain(config.Client.IPCPath)

			logs, err := blockchain.GetLogs(core.Contract{Hash: "x123"}, big.NewInt(4703824), nil)

			Expect(err).To(BeNil())
			Expect(len(logs)).To(Equal(0))

		})

	})

})
