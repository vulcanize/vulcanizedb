package integration

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/testing"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Reading contracts", func() {

	Describe("Getting a contract attribute", func() {
		It("retrieves the event log for a specific block and contract", func() {
			expectedLogZero := core.Log{
				BlockNumber: 4703824,
				TxHash:      "0xf896bfd1eb539d881a1a31102b78de9f25cd591bf1fe1924b86148c0b205fd5d",
				Address:     "0xd26114cd6ee289accf82350c8d8487fedb8a0c07",
				Topics: core.Topics{
					0: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					1: "0x000000000000000000000000fbb1b73c4f0bda4f67dca266ce6ef42f520fbb98",
					2: "0x000000000000000000000000d26114cd6ee289accf82350c8d8487fedb8a0c07",
				},
				Index: 19,
				Data:  "0x0000000000000000000000000000000000000000000000000c7d713b49da0000"}
			blockchain := geth.NewBlockChain(test_config.InfuraClient.IPCPath)
			contract := testing.SampleContract()

			logs, err := blockchain.GetLogs(contract, big.NewInt(4703824), nil)

			Expect(err).To(BeNil())
			Expect(len(logs)).To(Equal(3))
			Expect(logs[0]).To(Equal(expectedLogZero))

		})

		It("returns and empty log array when no events for a given block / contract combo", func() {
			blockchain := geth.NewBlockChain(test_config.InfuraClient.IPCPath)

			logs, err := blockchain.GetLogs(core.Contract{Hash: "x123"}, big.NewInt(4703824), nil)

			Expect(err).To(BeNil())
			Expect(len(logs)).To(Equal(0))

		})

	})

	Describe("Fetching Contract data", func() {
		It("returns the correct attribute for a real contract", func() {
			blockchain := geth.NewBlockChain(test_config.InfuraClient.IPCPath)

			contract := testing.SampleContract()
			var balance = new(big.Int)
			args := common.HexToHash("0xd26114cd6ee289accf82350c8d8487fedb8a0c07")
			err := blockchain.FetchContractData(contract.Abi, "0xd26114cd6ee289accf82350c8d8487fedb8a0c07", "balanceOf", args, &balance, 5167471)
			Expect(err).NotTo(HaveOccurred())
			expected := new(big.Int)
			expected.SetString("10897295492887612977137", 10)
			Expect(balance).To(Equal(expected))
		})
	})

})
