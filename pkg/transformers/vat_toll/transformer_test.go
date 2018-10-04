package vat_toll_test

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
	vat_toll_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_toll"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
)

var _ = Describe("Vat toll transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &vat_toll_mocks.MockVatTollRepository{}
		transformer := vat_toll.VatTollTransformer{
			Config:     vat_toll.VatTollConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vat_toll_mocks.MockVatTollConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_toll.VatTollConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_toll.VatTollConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &vat_toll_mocks.MockVatTollRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vat_toll_mocks.MockVatTollConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &vat_toll_mocks.MockVatTollRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    fetcher,
			Converter:  &vat_toll_mocks.MockVatTollConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{vat_toll.VatTollConfig.ContractAddresses, vat_toll.VatTollConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatTollSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &vat_toll_mocks.MockVatTollRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    fetcher,
			Converter:  &vat_toll_mocks.MockVatTollConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		mockConverter := &vat_toll_mocks.MockVatTollConverter{}
		mockRepository := &vat_toll_mocks.MockVatTollRepository{}
		headerID := int64(123)
		mockRepository.SetMissingHeaders([]core.Header{{Id: headerID}})
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := vat_toll.VatTollTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &vat_toll_mocks.MockVatTollConverter{}
		mockRepository := &vat_toll_mocks.MockVatTollRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := vat_toll.VatTollTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &vat_toll_mocks.MockVatTollConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatTollLog})
		repository := &vat_toll_mocks.MockVatTollRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthVatTollLog}))
	})

	It("returns error if converter returns error", func() {
		converter := &vat_toll_mocks.MockVatTollConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatTollLog})
		repository := &vat_toll_mocks.MockVatTollRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat toll model", func() {
		converter := &vat_toll_mocks.MockVatTollConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatTollLog})
		repository := &vat_toll_mocks.MockVatTollRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]vat_toll.VatTollModel{test_data.VatTollModel}))
	})

	It("returns error if repository returns error for create", func() {
		converter := &vat_toll_mocks.MockVatTollConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatTollLog})
		repository := &vat_toll_mocks.MockVatTollRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := vat_toll.VatTollTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
