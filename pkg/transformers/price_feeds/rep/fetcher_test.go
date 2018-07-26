package rep_test

import (
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
		mockBlockChain.SetGetLogsReturnLogs([]core.Log{{}})
		fetcher := rep.NewRepFetcher(mockBlockChain)
		header := core.Header{
			BlockNumber: 100,
			Hash:        "",
			Raw:         nil,
		}

		_, err := fetcher.FetchRepValue(header)

		Expect(err).NotTo(HaveOccurred())
		mockBlockChain.AssertGetLogsCalledWith(price_feeds.RepAddress, price_feeds.RepLogTopic0, big.NewInt(header.BlockNumber), big.NewInt(header.BlockNumber))
	})

	It("returns error if getting logs fails", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetLogsReturnLogs([]core.Log{{}})
		mockBlockChain.SetGetLogsReturnErr(fakes.FakeError)
		fetcher := rep.NewRepFetcher(mockBlockChain)

		_, err := fetcher.FetchRepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("returns no matching logs error if no logs returned", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		fetcher := rep.NewRepFetcher(mockBlockChain)

		_, err := fetcher.FetchRepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(price_feeds.ErrNoMatchingLog))
	})

	It("returns error if more than one matching logs returned", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetLogsReturnLogs([]core.Log{{}, {}})
		fetcher := rep.NewRepFetcher(mockBlockChain)

		_, err := fetcher.FetchRepValue(core.Header{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(price_feeds.ErrMultipleLogs))
	})
})
