// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package integration

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/geth/testing"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Reading contracts", func() {

	Describe("Getting a contract attribute", func() {
		It("retrieves the event log for a specific block and contract", func() {
			expectedLogZero := core.FullSyncLog{
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
			rawRpcClient, err := rpc.Dial(test_config.InfuraClient.IPCPath)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, test_config.InfuraClient.IPCPath)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			contract := testing.SampleContract()

			logs, err := blockChain.GetFullSyncLogs(contract, big.NewInt(4703824), nil)

			Expect(err).To(BeNil())
			Expect(len(logs)).To(Equal(3))
			Expect(logs[0]).To(Equal(expectedLogZero))
		})

		It("returns and empty log array when no events for a given block / contract combo", func() {
			rawRpcClient, err := rpc.Dial(test_config.InfuraClient.IPCPath)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, test_config.InfuraClient.IPCPath)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

			logs, err := blockChain.GetFullSyncLogs(core.Contract{Hash: "0x123"}, big.NewInt(4703824), nil)

			Expect(err).To(BeNil())
			Expect(len(logs)).To(Equal(0))
		})
	})

	Describe("Fetching Contract data", func() {
		It("returns the correct attribute for a real contract", func() {
			rawRpcClient, err := rpc.Dial(test_config.InfuraClient.IPCPath)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, test_config.InfuraClient.IPCPath)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

			contract := testing.SampleContract()
			var balance = new(big.Int)

			args := make([]interface{}, 1)
			args[0] = common.HexToHash("0xd26114cd6ee289accf82350c8d8487fedb8a0c07")

			err = blockChain.FetchContractData(contract.Abi, "0xd26114cd6ee289accf82350c8d8487fedb8a0c07", "balanceOf", args, &balance, 5167471)
			Expect(err).NotTo(HaveOccurred())
			expected := new(big.Int)
			expected.SetString("10897295492887612977137", 10)
			Expect(balance).To(Equal(expected))
		})
	})
})
