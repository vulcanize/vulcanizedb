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
	var mockClient *fakes.MockEthClient
	var blockChain *geth.BlockChain

	BeforeEach(func() {
		mockClient = fakes.NewMockEthClient()
		node := vulcCore.Node{}
		blockChain = geth.NewBlockChain(mockClient, node, cold_db.NewColdDbTransactionConverter())
	})

	Describe("getting a block", func() {
		It("fetches block from client", func() {
			mockClient.SetBlockByNumberReturnBlock(types.NewBlockWithHeader(&types.Header{}))
			blockNumber := int64(100)

			_, err := blockChain.GetBlockByNumber(blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertBlockByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
		})

		It("returns err if client returns err", func() {
			mockClient.SetBlockByNumberErr(fakes.FakeError)

			_, err := blockChain.GetBlockByNumber(100)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting a header", func() {
		It("fetches header from client", func() {
			blockNumber := int64(100)
			mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})

			_, err := blockChain.GetHeaderByNumber(blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockClient.AssertHeaderByNumberCalledWith(context.Background(), big.NewInt(blockNumber))
		})

		It("returns err if client returns err", func() {
			mockClient.SetHeaderByNumberErr(fakes.FakeError)

			_, err := blockChain.GetHeaderByNumber(100)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting logs with default FilterQuery", func() {
		It("fetches logs from client", func() {
			mockClient.SetFilterLogsReturnLogs([]types.Log{{}})
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
			mockClient.SetFilterLogsErr(fakes.FakeError)
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
			startingBlockNumber := big.NewInt(1)
			endingBlockNumber := big.NewInt(2)

			_, err := blockChain.GetLogs(contract, startingBlockNumber, endingBlockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("getting logs with a custom FilterQuery", func() {
		It("fetches logs from client", func() {
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

		It("returns err if client returns err", func() {
			mockClient.SetFilterLogsErr(fakes.FakeError)
			contract := vulcCore.Contract{Hash: common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex()}
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

	Describe("getting the most recent block number", func() {
		It("fetches latest header from client", func() {
			blockNumber := int64(100)
			mockClient.SetHeaderByNumberReturnHeader(&types.Header{Number: big.NewInt(blockNumber)})

			result := blockChain.LastBlock()

			mockClient.AssertHeaderByNumberCalledWith(context.Background(), nil)
			Expect(result).To(Equal(big.NewInt(blockNumber)))
		})
	})
})
