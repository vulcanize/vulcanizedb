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

package eth_test

import (
	"context"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Geth blockchain", func() {
	var (
		mockClient               *fakes.MockEthClient
		blockChain               *eth.BlockChain
		mockRpcClient            *fakes.MockRpcClient
		mockTransactionConverter *fakes.MockTransactionConverter
		node                     core.Node
	)

	BeforeEach(func() {
		mockClient = fakes.NewMockEthClient()
		mockRpcClient = fakes.NewMockRpcClient()
		mockTransactionConverter = fakes.NewMockTransactionConverter()
		node = core.Node{}
		blockChain = eth.NewBlockChain(mockClient, mockRpcClient, node, mockTransactionConverter)
	})

	Describe("getting a header", func() {
		Describe("default/mainnet", func() {
			It("fetches header from ethClient", func() {
				blockNumber := int64(100)
				mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})

				_, err := blockChain.GetHeaderByNumber(blockNumber)

				Expect(err).NotTo(HaveOccurred())
				mockClient.AssertHeaderByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
			})

			It("returns err if ethClient returns err", func() {
				mockClient.SetHeaderByNumberErr(fakes.FakeError)

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("fetches headers with multiple blocks", func() {
				_, err := blockChain.GetHeadersByNumbers([]int64{100, 99})

				Expect(err).NotTo(HaveOccurred())
				mockRpcClient.AssertBatchCalledWith("eth_getBlockByNumber", 2)
			})
		})

		Describe("POA/Kovan", func() {
			It("fetches header from rpcClient", func() {
				node.NetworkID = core.KOVAN_NETWORK_ID
				blockNumber := hexutil.Big(*big.NewInt(100))
				mockRpcClient.SetReturnPOAHeader(core.POAHeader{Number: &blockNumber})
				blockChain = eth.NewBlockChain(mockClient, mockRpcClient, node, fakes.NewMockTransactionConverter())

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).NotTo(HaveOccurred())
				mockRpcClient.AssertCallContextCalledWith(context.Background(), &core.POAHeader{}, "eth_getBlockByNumber")
			})

			It("returns err if rpcClient returns err", func() {
				node.NetworkID = core.KOVAN_NETWORK_ID
				mockRpcClient.SetCallContextErr(fakes.FakeError)
				blockChain = eth.NewBlockChain(mockClient, mockRpcClient, node, fakes.NewMockTransactionConverter())

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("returns error if returned header is empty", func() {
				node.NetworkID = core.KOVAN_NETWORK_ID
				blockChain = eth.NewBlockChain(mockClient, mockRpcClient, node, fakes.NewMockTransactionConverter())

				_, err := blockChain.GetHeaderByNumber(100)

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(eth.ErrEmptyHeader))
			})

			It("returns multiple headers with multiple blocknumbers", func() {
				node.NetworkID = core.KOVAN_NETWORK_ID
				blockNumber := hexutil.Big(*big.NewInt(100))
				mockRpcClient.SetReturnPOAHeaders([]core.POAHeader{{Number: &blockNumber}})

				_, err := blockChain.GetHeadersByNumbers([]int64{100, 99})

				Expect(err).NotTo(HaveOccurred())
				mockRpcClient.AssertBatchCalledWith("eth_getBlockByNumber", 2)
			})
		})
	})

	Describe("getting logs with a custom FilterQuery", func() {
		It("fetches logs from ethClient", func() {
			mockClient.SetFilterLogsReturnLogs([]types.Log{{}})
			address := common.HexToAddress("0x")
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)
			topic := common.HexToHash("0x")
			query := ethereum.FilterQuery{
				FromBlock: startingBlockNumber,
				ToBlock:   endingBlockNumber,
				Addresses: []common.Address{address},
				Topics:    [][]common.Hash{{topic}},
			}

			_, err := blockChain.GetEthLogsWithCustomQuery(query)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertFilterLogsCalledWith(context.Background(), query)
		})

		It("returns err if ethClient returns err", func() {
			mockClient.SetFilterLogsErr(fakes.FakeError)
			contract := core.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)
			query := ethereum.FilterQuery{
				FromBlock: startingBlockNumber,
				ToBlock:   endingBlockNumber,
				Addresses: []common.Address{common.HexToAddress(contract.Hash)},
				Topics:    nil,
			}

			_, err := blockChain.GetEthLogsWithCustomQuery(query)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting transactions", func() {
		It("fetches transaction for each hash", func() {
			_, err := blockChain.GetTransactions([]common.Hash{{}, {}})

			Expect(err).NotTo(HaveOccurred())
			mockRpcClient.AssertBatchCalledWith("eth_getTransactionByHash", 2)
		})

		It("converts rpc transaction to models", func() {
			_, err := blockChain.GetTransactions([]common.Hash{{}, {}})

			Expect(err).NotTo(HaveOccurred())
			Expect(mockTransactionConverter.ConvertRpcTransactionsToModelsCalled).To(BeTrue())
		})
	})

	Describe("getting the most recent block number", func() {
		It("fetches latest header from ethClient", func() {
			blockNumber := int64(100)
			mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})

			result, err := blockChain.LastBlock()
			Expect(err).NotTo(HaveOccurred())

			mockClient.AssertHeaderByNumberCalledWith(context.Background(), nil)
			Expect(result).To(Equal(big.NewInt(blockNumber)))
		})
	})

	Describe("getting storage at the given block", func() {
		var (
			account     = fakes.FakeAddress
			key         = fakes.FakeHash
			blockNumber = big.NewInt(rand.Int63())
		)

		It("fetches the storage at the given key, contract and block", func() {
			_, err := blockChain.GetStorageAt(account, key, blockNumber)
			Expect(err).NotTo(HaveOccurred())

			mockClient.AssertStorageAtCalledWith(context.Background(), account, key, blockNumber)
		})

		It("returns an error if the call to the eth client fails", func() {
			mockClient.SetStorageAtError(fakes.FakeError)

			_, err := blockChain.GetStorageAt(account, key, blockNumber)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(fakes.FakeError))
		})
	})
})
