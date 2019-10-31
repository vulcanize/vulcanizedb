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

package integration_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/eth/node"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Reading from the Geth blockchain", func() {
	var blockChain *eth.BlockChain

	BeforeEach(func() {
		rawRpcClient, err := rpc.Dial(test_config.TestClient.IPCPath)
		Expect(err).NotTo(HaveOccurred())
		rpcClient := client.NewRpcClient(rawRpcClient, test_config.TestClient.IPCPath)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		node := node.MakeNode(rpcClient)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		blockChain = eth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
	})

	It("reads two blocks", func(done Done) {
		blocks := fakes.NewMockBlockRepository()
		lastBlock, err := blockChain.LastBlock()
		Expect(err).NotTo(HaveOccurred())

		queriedBlocks := []int64{lastBlock.Int64() - 5, lastBlock.Int64() - 6}
		_, err = history.RetrieveAndUpdateBlocks(blockChain, blocks, queriedBlocks)
		Expect(err).NotTo(HaveOccurred())

		blocks.AssertCreateOrUpdateBlocksCallCountAndBlockNumbersEquals(2, []int64{lastBlock.Int64() - 5, lastBlock.Int64() - 6})
		close(done)
	}, 30)

	It("retrieves the genesis block and first block", func(done Done) {
		genesisBlock, err := blockChain.GetBlockByNumber(int64(0))
		Expect(err).ToNot(HaveOccurred())
		firstBlock, err := blockChain.GetBlockByNumber(int64(1))
		Expect(err).ToNot(HaveOccurred())
		lastBlockNumber, err := blockChain.LastBlock()

		Expect(err).NotTo(HaveOccurred())
		Expect(genesisBlock.Number).To(Equal(int64(0)))
		Expect(firstBlock.Number).To(Equal(int64(1)))
		Expect(lastBlockNumber.Int64()).To(BeNumerically(">", 0))
		close(done)
	}, 15)

	It("retrieves the node info", func(done Done) {
		node := blockChain.Node()
		mainnetID := float64(1)

		Expect(node.GenesisBlock).ToNot(BeNil())
		Expect(node.NetworkID).To(Equal(mainnetID))
		Expect(len(node.ID)).ToNot(BeZero())
		Expect(node.ClientName).ToNot(BeZero())

		close(done)
	}, 15)

	It("retrieves transaction", func() {
		// actual transaction: https://etherscan.io/tx/0x44d462f2a19ad267e276b234a62c542fc91c974d2e4754a325ca405f95440255
		txHash := common.HexToHash("0x44d462f2a19ad267e276b234a62c542fc91c974d2e4754a325ca405f95440255")
		transactions, err := blockChain.GetTransactions([]common.Hash{txHash})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(transactions)).To(Equal(1))
		expectedData := []byte{149, 227, 197, 11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			1, 160, 85, 105, 13, 157, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7,
			241, 202, 218, 90, 30, 178, 234, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 92, 155, 193, 43}
		expectedRaw := []byte{248, 201, 9, 132, 59, 154, 202, 0, 131, 1, 102, 93, 148, 44, 75, 208, 100, 185, 152, 131,
			128, 118, 250, 52, 26, 131, 208, 7, 252, 47, 165, 9, 87, 128, 184, 100, 149, 227, 197, 11, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 160, 85, 105, 13, 157, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 241, 202, 218, 90, 30, 178, 234, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 92, 155, 193, 43, 37, 160, 237, 184, 236, 248, 23, 152,
			53, 238, 44, 215, 181, 234, 229, 157, 246, 212, 178, 88, 25, 116, 134, 163, 124, 64, 2, 66, 25, 118, 1, 253, 27,
			101, 160, 36, 226, 116, 43, 147, 236, 124, 76, 227, 250, 228, 168, 22, 19, 248, 155, 248, 151, 219, 14, 1, 186,
			159, 35, 154, 22, 222, 123, 254, 147, 63, 221}
		expectedModel := core.TransactionModel{
			Data:     expectedData,
			From:     "0x3b08b99441086edd66f36f9f9aee733280698378",
			GasLimit: 91741,
			GasPrice: 1000000000,
			Hash:     "0x44d462f2a19ad267e276b234a62c542fc91c974d2e4754a325ca405f95440255",
			Nonce:    9,
			Raw:      expectedRaw,
			Receipt:  core.Receipt{},
			To:       "0x2c4bd064b998838076fa341a83d007fc2fa50957",
			TxIndex:  30,
			Value:    "0",
		}
		Expect(transactions[0]).To(Equal(expectedModel))
	})

	//Benchmarking test: remove skip to test performance of block retrieval
	XMeasure("retrieving n blocks", func(b Benchmarker) {
		b.Time("runtime", func() {
			var blocks []core.Block
			n := 10
			for i := 5327459; i > 5327459-n; i-- {
				block, err := blockChain.GetBlockByNumber(int64(i))
				Expect(err).ToNot(HaveOccurred())
				blocks = append(blocks, block)
			}
			Expect(len(blocks)).To(Equal(n))
		})
	}, 10)
})
