package vat_grab_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
)

var _ = Describe("Vat grab transformer", func() {
	var (
		config      = vat_grab.VatGrabConfig
		converter   mocks.MockLogNoteConverter
		repository  mocks.MockRepository
		fetcher     mocks.MockLogFetcher
		transformer shared.Transformer
	)

	BeforeEach(func() {
		repository = mocks.MockRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockLogNoteConverter{}
		transformer = factories.LogNoteTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}.NewLogNoteTransformer(nil, nil)
	})

	It("sets the blockchain and database", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers for block numbers specified in config", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_grab.VatGrabConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_grab.VatGrabConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{vat_grab.VatGrabConfig.ContractAddresses, vat_grab.VatGrabConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatGrabSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		headerID := int64(123)
		repository.SetMissingHeaders([]core.Header{{Id: headerID}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatGrabLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthVatGrabLog}))
	})

	It("returns error if converter returns error", func() {
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatGrabLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat grab model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatGrabLog})
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		converter.SetReturnModels([]interface{}{test_data.VatGrabModel})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels[0]).To(Equal(test_data.VatGrabModel))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatGrabLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
