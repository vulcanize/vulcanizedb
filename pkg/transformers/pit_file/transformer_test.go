package pit_file_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	pit_file_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/pit_file"
)

var _ = Describe("Pit file transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &pit_file_mocks.MockPitFileRepository{}
		transformer := pit_file.PitFileTransformer{
			Config:     pit_file.PitFileConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &pit_file_mocks.MockPitFileConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(pit_file.PitFileConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(pit_file.PitFileConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &pit_file_mocks.MockPitFileRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := pit_file.PitFileTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &pit_file_mocks.MockPitFileConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &pit_file_mocks.MockPitFileRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := pit_file.PitFileTransformer{
			Fetcher:    fetcher,
			Converter:  &pit_file_mocks.MockPitFileConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddress).To(Equal(pit_file.PitFileConfig.ContractAddress))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.PitFileSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &pit_file_mocks.MockPitFileRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := pit_file.PitFileTransformer{
			Fetcher:    fetcher,
			Converter:  &pit_file_mocks.MockPitFileConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &pit_file_mocks.MockPitFileConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileLog})
		repository := &pit_file_mocks.MockPitFileRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := pit_file.PitFileTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedContractAddress).To(Equal(pit_file.PitFileConfig.ContractAddress))
		Expect(converter.PassedContractABI).To(Equal(pit_file.PitFileConfig.ContractAbi))
		Expect(converter.PassedLog).To(Equal(test_data.EthPitFileLog))
	})

	It("returns error if converter returns error", func() {
		converter := &pit_file_mocks.MockPitFileConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileLog})
		repository := &pit_file_mocks.MockPitFileRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := pit_file.PitFileTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists pit file model", func() {
		converter := &pit_file_mocks.MockPitFileConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileLog})
		repository := &pit_file_mocks.MockPitFileRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := pit_file.PitFileTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedTransactionIndex).To(Equal(test_data.EthPitFileLog.TxIndex))
		Expect(repository.PassedModel).To(Equal(test_data.PitFileModel))
	})

	It("returns error if repository returns error for create", func() {
		converter := &pit_file_mocks.MockPitFileConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileLog})
		repository := &pit_file_mocks.MockPitFileRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := pit_file.PitFileTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
