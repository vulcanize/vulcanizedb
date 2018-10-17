package vat_flux_test

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
	vat_flux_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_flux"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
)

var _ = Describe("Vat flux transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		transformer := vat_flux.VatFluxTransformer{
			Config:     vat_flux.VatFluxConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vat_flux_mocks.MockVatFlux{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_flux.VatFluxConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_flux.VatFluxConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vat_flux_mocks.MockVatFlux{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    fetcher,
			Converter:  &vat_flux_mocks.MockVatFlux{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{vat_flux.VatFluxConfig.ContractAddresses, vat_flux.VatFluxConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatFluxSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    fetcher,
			Converter:  &vat_flux_mocks.MockVatFlux{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &vat_flux_mocks.MockVatFlux{}
		mockRepository := &vat_flux_mocks.MockVatFluxRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := vat_flux.VatFluxTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &vat_flux_mocks.MockVatFlux{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.VatFluxLog}))
	})

	It("returns error if converter returns error", func() {
		converter := &vat_flux_mocks.MockVatFlux{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat flux model", func() {
		converter := &vat_flux_mocks.MockVatFlux{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]vat_flux.VatFluxModel{test_data.VatFluxModel}))
	})

	It("returns error if repository returns error for create", func() {
		converter := &vat_flux_mocks.MockVatFlux{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository := &vat_flux_mocks.MockVatFluxRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := vat_flux.VatFluxTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
