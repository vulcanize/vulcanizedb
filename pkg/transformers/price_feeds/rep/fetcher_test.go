package rep_test

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/rep"
	"math/big"
)

var _ = Describe("Rep fetcher", func() {
	It("gets logs describing updated rep/usd value", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{}})
		contractAddress := "rep-contract-address"
		fetcher := rep.NewRepFetcher(mockBlockChain, contractAddress)
		blockNumber := int64(100)
		header := core.Header{
			BlockNumber: blockNumber,
			Hash:        "",
			Raw:         nil,
		}

		_, err := fetcher.FetchRepValue(header)

		Expect(err).NotTo(HaveOccurred())
		expectedQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(blockNumber),
			ToBlock:   big.NewInt(blockNumber),
			Addresses: []common.Address{common.HexToAddress(contractAddress)},
			Topics:    [][]common.Hash{{common.HexToHash(price_feeds.RepLogTopic0)}},
		}
		mockBlockChain.AssertGetEthLogsWithCustomQueryCalledWith(expectedQuery)
	})

	It("returns error if getting logs fails", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryErr(fakes.FakeError)
		fetcher := rep.NewRepFetcher(mockBlockChain, "rep-contract-address")

		_, err := fetcher.FetchRepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("returns no matching logs error if no logs returned", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		fetcher := rep.NewRepFetcher(mockBlockChain, "rep-contract-address")

		_, err := fetcher.FetchRepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(price_feeds.ErrNoMatchingLog))
	})

	It("returns error if more than one matching logs returned", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{}, {}})
		fetcher := rep.NewRepFetcher(mockBlockChain, "rep-contract-address")

		_, err := fetcher.FetchRepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(price_feeds.ErrMultipleLogs))
	})
})
