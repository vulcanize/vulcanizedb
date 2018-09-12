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

package flip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	flip_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/cat_file/flip"
)

var _ = Describe("Cat file flip transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &flip_mocks.MockCatFileFlipRepository{}
		transformer := flip.CatFileFlipTransformer{
			Config:     cat_file.CatFileConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &flip_mocks.MockCatFileFlipConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(cat_file.CatFileConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(cat_file.CatFileConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &flip_mocks.MockCatFileFlipRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &flip_mocks.MockCatFileFlipConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &flip_mocks.MockCatFileFlipRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    fetcher,
			Converter:  &flip_mocks.MockCatFileFlipConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{cat_file.CatFileConfig.ContractAddresses, cat_file.CatFileConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.CatFileFlipSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &flip_mocks.MockCatFileFlipRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    fetcher,
			Converter:  &flip_mocks.MockCatFileFlipConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		mockConverter := &flip_mocks.MockCatFileFlipConverter{}
		mockRepository := &flip_mocks.MockCatFileFlipRepository{}
		headerID := int64(123)
		mockRepository.SetMissingHeaders([]core.Header{{Id: headerID}})
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := flip.CatFileFlipTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &flip_mocks.MockCatFileFlipConverter{}
		mockRepository := &flip_mocks.MockCatFileFlipRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := flip.CatFileFlipTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &flip_mocks.MockCatFileFlipConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthCatFileFlipLog})
		repository := &flip_mocks.MockCatFileFlipRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthCatFileFlipLog}))
	})

	It("returns error if converter returns error", func() {
		converter := &flip_mocks.MockCatFileFlipConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthCatFileFlipLog})
		repository := &flip_mocks.MockCatFileFlipRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists cat file flip model", func() {
		converter := &flip_mocks.MockCatFileFlipConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthCatFileFlipLog})
		repository := &flip_mocks.MockCatFileFlipRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]flip.CatFileFlipModel{test_data.CatFileFlipModel}))
	})

	It("returns error if repository returns error for create", func() {
		converter := &flip_mocks.MockCatFileFlipConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthCatFileFlipLog})
		repository := &flip_mocks.MockCatFileFlipRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := flip.CatFileFlipTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
