package vat_slip_test

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
	vat_slip_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_slip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
)

var _ = Describe("Vat slip transformer", func() {
	var (
		config      shared.TransformerConfig
		converter   *vat_slip_mocks.MockVatSlipConverter
		fetcher     *mocks.MockLogFetcher
		repository  *vat_slip_mocks.MockVatSlipRepository
		transformer vat_slip.VatSlipTransformer
	)

	BeforeEach(func() {
		config = vat_slip.VatSlipConfig
		converter = &vat_slip_mocks.MockVatSlipConverter{}
		fetcher = &mocks.MockLogFetcher{}
		repository = &vat_slip_mocks.MockVatSlipRepository{}
		transformer = vat_slip.VatSlipTransformer{
			Config:     config,
			Converter:  converter,
			Fetcher:    fetcher,
			Repository: repository,
		}
	})

	It("gets missing headers for block numbers specified in config", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_slip.VatSlipConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_slip.VatSlipConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersErr(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		headerOne := core.Header{BlockNumber: GinkgoRandomSeed()}
		headerTwo := core.Header{BlockNumber: GinkgoRandomSeed()}
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{vat_slip.VatSlipConfig.ContractAddresses, vat_slip.VatSlipConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatSlipSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: GinkgoRandomSeed()}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		headerID := GinkgoRandomSeed()
		repository.SetMissingHeaders([]core.Header{{Id: headerID}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{{Id: GinkgoRandomSeed()}})
		repository.SetMarkHeaderCheckedErr(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatSlipLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: GinkgoRandomSeed()}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthVatSlipLog}))
	})

	It("returns error if converter returns error", func() {
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatSlipLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: GinkgoRandomSeed()}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat slip model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatSlipLog})
		fakeHeader := core.Header{BlockNumber: GinkgoRandomSeed(), Id: GinkgoRandomSeed()}
		repository.SetMissingHeaders([]core.Header{fakeHeader})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]vat_slip.VatSlipModel{test_data.VatSlipModel}))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatSlipLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: GinkgoRandomSeed(), Id: GinkgoRandomSeed()}})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
