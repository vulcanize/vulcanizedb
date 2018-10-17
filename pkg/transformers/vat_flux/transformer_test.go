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

type setupOptions struct {
	setMissingHeadersError bool
	setFetcherError        bool
	setConverterError      bool
	setCreateError         bool
	fetchedLogs            []types.Log
	missingHeaders         []core.Header
}

var _ = Describe("Vat flux transformer", func() {
	var (
		config      shared.TransformerConfig
		converter   *vat_flux_mocks.MockVatFlux
		fetcher     *mocks.MockLogFetcher
		repository  *vat_flux_mocks.MockVatFluxRepository
		transformer vat_flux.VatFluxTransformer
	)

	BeforeEach(func() {
		config = vat_flux.VatFluxConfig
		converter = &vat_flux_mocks.MockVatFlux{}
		fetcher = &mocks.MockLogFetcher{}
		repository = &vat_flux_mocks.MockVatFluxRepository{}
		transformer = vat_flux.VatFluxTransformer{
			Config:     config,
			Converter:  converter,
			Fetcher:    fetcher,
			Repository: repository,
		}
	})

	It("gets missing headers for block numbers specified in config", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_flux.VatFluxConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_flux.VatFluxConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersErr(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks the header as checked when there are no logs", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		repository.SetMissingHeaders([]core.Header{header})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.MarkHeaderCheckedPassedHeaderID).To(Equal(header.Id))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{vat_flux.VatFluxConfig.ContractAddresses, vat_flux.VatFluxConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatFluxSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		repository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := vat_flux.VatFluxTransformer{
			Converter:  converter,
			Fetcher:    mockFetcher,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
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
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat flux model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]vat_flux.VatFluxModel{test_data.VatFluxModel}))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
