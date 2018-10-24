package vat_flux_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	vat_flux_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_flux"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
	"math/rand"
)

var _ = Describe("Vat flux transformer", func() {
	var (
		config      = vat_flux.VatFluxConfig
		fetcher     mocks.MockLogFetcher
		converter   vat_flux_mocks.MockVatFluxConverter
		repository  vat_flux_mocks.MockVatFluxRepository
		transformer shared.Transformer
		headerOne   core.Header
		headerTwo   core.Header
	)

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = vat_flux_mocks.MockVatFluxConverter{}
		repository = vat_flux_mocks.MockVatFluxRepository{}
		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		headerTwo = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		transformer = factories.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}.NewTransformer(nil, nil)
	})

	It("sets the blockchain and database", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers for block numbers specified in config", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersErr(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{config.ContractAddresses, config.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatFluxSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks the header as checked when there are no logs", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.MarkHeaderCheckedPassedHeaderID).To(Equal(headerOne.Id))
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetMarkHeaderCheckedErr(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.VatFluxLog}))
	})

	It("returns error if converter returns error", func() {
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat flux model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels).To(Equal([]interface{}{test_data.VatFluxModel}))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.VatFluxLog})
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
