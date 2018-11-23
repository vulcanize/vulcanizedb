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

package every_block_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/generic/every_block"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var _ = Describe("every_block Getter", func() {
	blockNumber := int64(5502914)

	Describe("stopped", func() {
		It("checks whether or not the contract has been stopped", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetStoppedStatus(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult bool
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "stopped", nil, &expected, blockNumber)
		})

		It("gets dai token's stopped status at the given block height", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetStoppedStatus(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(false))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetStoppedStatus("", "", 0)

			Expect(result).To(Equal(false))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("stopped"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})

	Describe("owner", func() {
		It("checks what the contract's owner address is", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetOwner(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult common.Address
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "owner", nil, &expected, blockNumber)
		})

		It("gets dai token's owner address at the given block height", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetOwner(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := common.HexToAddress("0x0000000000000000000000000000000000000000")
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetOwner("", "", 0)

			expectedResult := new(common.Address)
			Expect(result).To(Equal(*expectedResult))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("owner"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})

	Describe("hash name", func() {
		It("checks the contract's name", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetHashName(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult common.Hash
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "name", nil, &expected, blockNumber)
		})

		It("gets dai token's name at the given blockheight", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetHashName(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := common.HexToHash("0x44616920537461626c65636f696e2076312e3000000000000000000000000000")
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetHashName("", "", 0)

			expectedResult := new(common.Hash)
			Expect(result).To(Equal(*expectedResult))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})

	Describe("hash symbol", func() {
		It("checks the contract's symbol", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetHashSymbol(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult common.Hash
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "symbol", nil, &expected, blockNumber)
		})

		It("gets dai token's symbol at the given blockheight", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetHashSymbol(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := common.HexToHash("0x4441490000000000000000000000000000000000000000000000000000000000")
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetHashSymbol("", "", 0)

			expectedResult := new(common.Hash)
			Expect(result).To(Equal(*expectedResult))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("symbol"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})

	Describe("decimals", func() {
		It("checks what the token's number of decimals", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetDecimals(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult big.Int
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "decimals", nil, &expected, blockNumber)
		})

		It("gets dai token's number of decimals at the given block height", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetDecimals(constants.DaiAbiString, constants.DaiContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := big.Int{}
			expectedResult.SetString("18", 10)
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetDecimals("", "", 0)

			expectedResult := new(big.Int)
			Expect(result).To(Equal(*expectedResult))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("decimals"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})

	Describe("string name", func() {
		It("checks the contract's name", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetStringName(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult string
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "name", nil, &expected, blockNumber)
		})

		It("gets tusd token's name at the given blockheight", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetStringName(constants.TusdAbiString, constants.TusdContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := "TrueUSD"
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetStringName("", "", 0)

			expectedResult := new(string)
			Expect(result).To(Equal(*expectedResult))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})

	Describe("string symbol", func() {
		It("checks the contract's symbol", func() {
			fakeBlockChain := fakes.NewMockBlockChain()
			testGetter := every_block.NewGetter(fakeBlockChain)
			testAbi := "testAbi"
			testContractAddress := "testContractAddress"

			_, err := testGetter.GetStringSymbol(testAbi, testContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			var expectedResult string
			expected := &expectedResult
			fakeBlockChain.AssertFetchContractDataCalledWith(testAbi, testContractAddress, "symbol", nil, &expected, blockNumber)
		})

		It("gets tusd token's symbol at the given blockheight", func() {
			infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
			rawRpcClient, err := rpc.Dial(infuraIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
			blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			realGetter := every_block.NewGetter(blockChain)
			result, err := realGetter.GetStringName(constants.TusdAbiString, constants.TusdContractAddress, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedResult := "TrueUSD"
			Expect(result).To(Equal(expectedResult))
		})

		It("returns an error if the call to the blockchain fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetFetchContractDataErr(fakes.FakeError)
			errorGetter := every_block.NewGetter(blockChain)
			result, err := errorGetter.GetStringSymbol("", "", 0)

			expectedResult := new(string)
			Expect(result).To(Equal(*expectedResult))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("symbol"))
			Expect(err.Error()).To(ContainSubstring(fakes.FakeError.Error()))
		})
	})
})
