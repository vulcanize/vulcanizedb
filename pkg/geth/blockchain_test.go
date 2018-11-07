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

package geth_test

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	vulcCore "github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/converters/cold_db"
)

var _ = Describe("Geth blockchain", func() {
	Describe("getting a block", func() {
		It("fetches block from client", func() {
			mockClient := fakes.NewMockEthClient()
			mockClient.SetBlockByNumberReturnBlock(types.NewBlockWithHeader(&types.Header{}))
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())
			blockNumber := int64(100)

			_, err := blockChain.GetBlockByNumber(blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertBlockByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
		})

		It("returns err if client returns err", func() {
			mockClient := fakes.NewMockEthClient()
			mockClient.SetBlockByNumberErr(fakes.FakeError)
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())

			_, err := blockChain.GetBlockByNumber(100)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting a header", func() {
		It("fetches header from client", func() {
			mockClient := fakes.NewMockEthClient()
			blockNumber := int64(100)
			mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())

			_, err := blockChain.GetHeaderByNumber(blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertHeaderByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
		})

		It("returns err if client returns err", func() {
			mockClient := fakes.NewMockEthClient()
			mockClient.SetHeaderByNumberErr(fakes.FakeError)
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())

			_, err := blockChain.GetHeaderByNumber(100)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting logs", func() {
		It("fetches logs from client", func() {
			mockClient := fakes.NewMockEthClient()
			mockClient.SetFilterLogsReturnLogs([]types.Log{{}})
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)

			_, err := blockChain.GetLogs(contract, startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedQuery := ethereum.FilterQuery{
				FromBlock: startingBlockNumber,
				ToBlock:   endingBlockNumber,
				Addresses: []common.Address{common.HexToAddress(contract.Hash)},
			}
			mockClient.AssertFilterLogsCalledWith(context.Background(), expectedQuery)
		})

		It("returns err if client returns err", func() {
			mockClient := fakes.NewMockEthClient()
			mockClient.SetFilterLogsErr(fakes.FakeError)
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)

			_, err := blockChain.GetLogs(contract, startingBlockNumber, endingBlockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting the most recent block number", func() {
		It("fetches latest header from client", func() {
			mockClient := fakes.NewMockEthClient()
			blockNumber := int64(100)
			mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})
			node := vulcCore.Node{}
			blockChain := geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())

			result := blockChain.LastBlock()

			mockClient.AssertHeaderByNumberCalledWith(context.Background(), nil)
			Expect(result).To(Equal(big.NewInt(blockNumber)))
		})
	})
})
