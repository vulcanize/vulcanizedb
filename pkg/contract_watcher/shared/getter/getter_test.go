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

package getter_test

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/shared/getter"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"
	rpc2 "github.com/makerdao/vulcanizedb/pkg/eth/converters/rpc"
	"github.com/makerdao/vulcanizedb/pkg/eth/node"
	"github.com/makerdao/vulcanizedb/test_config"
)

var _ = Describe("Interface Getter", func() {
	Describe("GetAbi", func() {
		It("Constructs and returns a custom abi based on results from supportsInterface calls", func() {
			expectedABI := `[` + constants.AddrChangeInterface + `,` + constants.NameChangeInterface + `,` + constants.ContentChangeInterface + `,` + constants.AbiChangeInterface + `,` + constants.PubkeyChangeInterface + `]`
			con := test_config.TestClient
			testIPC := con.IPCPath
			blockNumber := int64(6885696)
			rawRpcClient, err := rpc.Dial(testIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, testIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := eth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			interfaceGetter := getter.NewInterfaceGetter(blockChain)
			abi, err := interfaceGetter.GetABI(constants.PublicResolverAddress, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(abi).To(Equal(expectedABI))
			_, err = eth.ParseAbi(abi)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
