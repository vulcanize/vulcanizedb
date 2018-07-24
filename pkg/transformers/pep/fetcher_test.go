package pep_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pep"
	"math/big"
)

var _ = Describe("Pep fetcher", func() {
	It("calls contract to peek mkr/usd value", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{}})
		fetcher := pep.NewPepFetcher(mockBlockChain)
		blockNumber := int64(100)
		header := core.Header{
			BlockNumber: blockNumber,
			Hash:        "",
			Raw:         nil,
		}

		_, err := fetcher.FetchPepValue(header)

		Expect(err).NotTo(HaveOccurred())
		expectedQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(blockNumber),
			ToBlock:   big.NewInt(blockNumber),
			Addresses: []common.Address{common.HexToAddress(pep.PepAddress)},
			Topics:    [][]common.Hash{{common.HexToHash(pep.LogValueTopic0)}},
		}
		mockBlockChain.AssertGetEthLogsWithCustomQueryCalledWith(expectedQuery)
	})

	It("returns error if contract call fails", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{}})
		mockBlockChain.SetGetEthLogsWithCustomQueryErr(fakes.FakeError)
		fetcher := pep.NewPepFetcher(mockBlockChain)

		_, err := fetcher.FetchPepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("returns no matching logs error if no logs returned", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		fetcher := pep.NewPepFetcher(mockBlockChain)

		_, err := fetcher.FetchPepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(pep.ErrNoMatchingLog))
	})

	It("returns error if more than one matching logs returned", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{}, {}})
		fetcher := pep.NewPepFetcher(mockBlockChain)

		_, err := fetcher.FetchPepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(pep.ErrMultipleLogs))
	})
})
