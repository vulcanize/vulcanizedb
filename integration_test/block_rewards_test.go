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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/eth"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/eth/node"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Rewards calculations", func() {

	It("calculates a block reward for a real block", func() {
		rawRPCClient, err := rpc.Dial(test_config.TestClient.IPCPath)
		Expect(err).NotTo(HaveOccurred())
		rpcClient := client.NewRPCClient(rawRPCClient, test_config.TestClient.IPCPath)
		ethClient := ethclient.NewClient(rawRPCClient)
		blockChainClient := client.NewEthClient(ethClient)
		node := node.MakeNode(rpcClient)
		transactionConverter := vRpc.NewRPCTransactionConverter(ethClient)
		blockChain := eth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
		block, err := blockChain.GetBlockByNumber(1071819)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.Reward).To(Equal("5313550000000000000"))
	})

	It("calculates an uncle reward for a real block", func() {
		rawRPCClient, err := rpc.Dial(test_config.TestClient.IPCPath)
		Expect(err).NotTo(HaveOccurred())
		rpcClient := client.NewRPCClient(rawRPCClient, test_config.TestClient.IPCPath)
		ethClient := ethclient.NewClient(rawRPCClient)
		blockChainClient := client.NewEthClient(ethClient)
		node := node.MakeNode(rpcClient)
		transactionConverter := vRpc.NewRPCTransactionConverter(ethClient)
		blockChain := eth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
		block, err := blockChain.GetBlockByNumber(1071819)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.UnclesReward).To(Equal("6875000000000000000"))
	})

})
