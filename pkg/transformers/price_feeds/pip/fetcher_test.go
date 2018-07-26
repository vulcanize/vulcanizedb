package pip_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds/pip"
)

var _ = Describe("Pip fetcher", func() {
	It("returns error if fetching logs fails", func() {
		chain := fakes.NewMockBlockChain()
		chain.SetGetEthLogsWithCustomQueryErr(fakes.FakeError)
		fetcher := pip.NewPipFetcher(chain, "pip-contract-address")

		_, err := fetcher.FetchPipValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("returns no matching logs error if no logs returned", func() {
		chain := fakes.NewMockBlockChain()
		fetcher := pip.NewPipFetcher(chain, "pip-contract-address")

		_, err := fetcher.FetchPipValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(price_feeds.ErrNoMatchingLog))
	})

	Describe("when matching log found", func() {
		It("calls contract to peek current eth/usd value", func() {
			blockNumber := uint64(12345)
			chain := fakes.NewMockBlockChain()
			chain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{BlockNumber: blockNumber}})
			contractAddress := "pip-contract-address"
			fetcher := pip.NewPipFetcher(chain, contractAddress)

			_, err := fetcher.FetchPipValue(core.Header{})

			Expect(err).NotTo(HaveOccurred())
			chain.AssertFetchContractDataCalledWith(price_feeds.PipMedianizerABI, contractAddress, price_feeds.PeekMethodName, nil, &[]interface{}{[32]byte{}, false}, int64(blockNumber))
		})

		It("returns error if contract call fails", func() {
			chain := fakes.NewMockBlockChain()
			chain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{BlockNumber: uint64(12345)}})
			chain.SetFetchContractDataErr(fakes.FakeError)
			fetcher := pip.NewPipFetcher(chain, "pip-contract-address")

			_, err := fetcher.FetchPipValue(core.Header{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
