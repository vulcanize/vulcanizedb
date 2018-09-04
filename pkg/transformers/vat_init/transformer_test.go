package vat_init_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	vat_init_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_init"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
)

var _ = Describe("Vat init transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &vat_init_mocks.MockVatInitRepository{}
		transformer := vat_init.VatInitTransformer{
			Config:     vat_init.VatInitConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vat_init_mocks.MockVatInitConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_init.VatInitConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_init.VatInitConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &vat_init_mocks.MockVatInitRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := vat_init.VatInitTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vat_init_mocks.MockVatInitConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &vat_init_mocks.MockVatInitRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := vat_init.VatInitTransformer{
			Fetcher:    fetcher,
			Converter:  &vat_init_mocks.MockVatInitConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddress).To(Equal(vat_init.VatInitConfig.ContractAddress))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatInitSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &vat_init_mocks.MockVatInitRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_init.VatInitTransformer{
			Fetcher:    fetcher,
			Converter:  &vat_init_mocks.MockVatInitConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &vat_init_mocks.MockVatInitConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatInitLog})
		repository := &vat_init_mocks.MockVatInitRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_init.VatInitTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLog).To(Equal(test_data.EthVatInitLog))
	})

	It("returns error if converter returns error", func() {
		converter := &vat_init_mocks.MockVatInitConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatInitLog})
		repository := &vat_init_mocks.MockVatInitRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_init.VatInitTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat init model", func() {
		converter := &vat_init_mocks.MockVatInitConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatInitLog})
		repository := &vat_init_mocks.MockVatInitRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := vat_init.VatInitTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModel).To(Equal(test_data.VatInitModel))
	})

	It("returns error if repository returns error for create", func() {
		converter := &vat_init_mocks.MockVatInitConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatInitLog})
		repository := &vat_init_mocks.MockVatInitRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := vat_init.VatInitTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
