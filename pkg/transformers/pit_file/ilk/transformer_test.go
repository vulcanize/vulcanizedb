// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ilk_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	pit_file_ilk_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/pit_file/ilk"
)

var _ = Describe("Pit file ilk transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		transformer := ilk.PitFileIlkTransformer{
			Config:     pit_file.PitFileConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &pit_file_ilk_mocks.MockPitFileIlkConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(pit_file.PitFileConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(pit_file.PitFileConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &pit_file_ilk_mocks.MockPitFileIlkConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    fetcher,
			Converter:  &pit_file_ilk_mocks.MockPitFileIlkConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{pit_file.PitFileConfig.ContractAddresses, pit_file.PitFileConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.PitFileIlkSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    fetcher,
			Converter:  &pit_file_ilk_mocks.MockPitFileIlkConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		mockConverter := &pit_file_ilk_mocks.MockPitFileIlkConverter{}
		mockRepository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		headerID := int64(123)
		mockRepository.SetMissingHeaders([]core.Header{{Id: headerID}})
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := ilk.PitFileIlkTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &pit_file_ilk_mocks.MockPitFileIlkConverter{}
		mockRepository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := ilk.PitFileIlkTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &pit_file_ilk_mocks.MockPitFileIlkConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthPitFileIlkLog}))
	})

	It("returns error if converter returns error", func() {
		converter := &pit_file_ilk_mocks.MockPitFileIlkConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists pit file model", func() {
		converter := &pit_file_ilk_mocks.MockPitFileIlkConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]ilk.PitFileIlkModel{test_data.PitFileIlkModel}))
	})

	It("returns error if repository returns error for create", func() {
		converter := &pit_file_ilk_mocks.MockPitFileIlkConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository := &pit_file_ilk_mocks.MockPitFileIlkRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := ilk.PitFileIlkTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
