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

package vow_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	vow_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/drip_file/vow"
)

var _ = Describe("Drip file vow transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &vow_mocks.MockDripFileVowRepository{}
		transformer := vow.DripFileVowTransformer{
			Config:     drip_file.DripFileConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vow_mocks.MockDripFileVowConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(drip_file.DripFileConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(drip_file.DripFileConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &vow_mocks.MockDripFileVowRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := vow.DripFileVowTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &vow_mocks.MockDripFileVowConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &vow_mocks.MockDripFileVowRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := vow.DripFileVowTransformer{
			Fetcher:    fetcher,
			Converter:  &vow_mocks.MockDripFileVowConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{drip_file.DripFileConfig.ContractAddresses, drip_file.DripFileConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.DripFileVowSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &vow_mocks.MockDripFileVowRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vow.DripFileVowTransformer{
			Fetcher:    fetcher,
			Converter:  &vow_mocks.MockDripFileVowConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		mockConverter := &vow_mocks.MockDripFileVowConverter{}
		mockRepository := &vow_mocks.MockDripFileVowRepository{}
		headerID := int64(123)
		mockRepository.SetMissingHeaders([]core.Header{{Id: headerID}})
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := vow.DripFileVowTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &vow_mocks.MockDripFileVowConverter{}
		mockRepository := &vow_mocks.MockDripFileVowRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := vow.DripFileVowTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &vow_mocks.MockDripFileVowConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripFileVowLog})
		repository := &vow_mocks.MockDripFileVowRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vow.DripFileVowTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthDripFileVowLog}))
	})

	It("returns error if converter returns error", func() {
		converter := &vow_mocks.MockDripFileVowConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripFileVowLog})
		repository := &vow_mocks.MockDripFileVowRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := vow.DripFileVowTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists drip file model", func() {
		converter := &vow_mocks.MockDripFileVowConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripFileVowLog})
		repository := &vow_mocks.MockDripFileVowRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := vow.DripFileVowTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]vow.DripFileVowModel{test_data.DripFileVowModel}))
	})

	It("returns error if repository returns error for create", func() {
		converter := &vow_mocks.MockDripFileVowConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripFileVowLog})
		repository := &vow_mocks.MockDripFileVowRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := vow.DripFileVowTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
