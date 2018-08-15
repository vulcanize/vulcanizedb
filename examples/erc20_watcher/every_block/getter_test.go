// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package every_block_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/every_block"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

var _ = Describe("ERC20 Getter", func() {
	blockNumber := int64(5502914)

	Describe("totalSupply", func() {
		It("gets total supply data from the blockchain with the correct arguments", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetTotalSupply(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := big.Int{}
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "totalSupply", nil, &expected, blockNumber)
		})

		It("gets dai token's total supply at the given block height", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetTotalSupply(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := big.Int{}
			expectedResult.SetString("27647235749155415536952630", 10)
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetTotalSupply("", "", 0)

			Expect(result.String()).To(Equal("0"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("totalSupply"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})
})
