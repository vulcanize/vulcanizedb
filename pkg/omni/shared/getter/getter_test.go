// VulcanizeDB
// Copyright © 2018 Vulcanize

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

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/getter"
)

var _ = Describe("Interface Getter", func() {
	Describe("GetAbi", func() {
		It("Constructs and returns a custom abi based on results from supportsInterface calls", func() {
			expectedABI := `[` + constants.AddrChangeInterface + `,` + constants.NameChangeInterface + `,` + constants.ContentChangeInterface + `,` + constants.AbiChangeInterface + `,` + constants.PubkeyChangeInterface + `]`

			blockNumber := int64(6885696)
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			interfaceGetter := getter.NewInterfaceGetter(blockChain)
			abi := interfaceGetter.GetABI(constants.PublicResolverAddress, blockNumber)
			Expect(abi).To(Equal(expectedABI))
			_, err = geth.ParseAbi(abi)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
