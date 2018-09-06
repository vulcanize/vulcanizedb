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

package debt_ceiling_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	debt_ceiling_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/pit_file/debt_ceiling"
)

var _ = Describe("", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Config:     pit_file.PitFileConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(pit_file.PitFileConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(pit_file.PitFileConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    fetcher,
			Converter:  &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddress).To(Equal(pit_file.PitFileConfig.ContractAddress))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.PitFileDebtCeilingSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    fetcher,
			Converter:  &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileDebtCeilingLog})
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedContractAddress).To(Equal(pit_file.PitFileConfig.ContractAddress))
		Expect(converter.PassedContractABI).To(Equal(pit_file.PitFileConfig.ContractAbi))
		Expect(converter.PassedLog).To(Equal(test_data.EthPitFileDebtCeilingLog))
	})

	It("returns error if converter returns error", func() {
		converter := &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileDebtCeilingLog})
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists pit file model", func() {
		converter := &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileDebtCeilingLog})
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModel).To(Equal(test_data.PitFileDebtCeilingModel))
	})

	It("returns error if repository returns error for create", func() {
		converter := &debt_ceiling_mocks.MockPitFileDebtCeilingConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileDebtCeilingLog})
		repository := &debt_ceiling_mocks.MockPitFileDebtCeilingRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := debt_ceiling.PitFileDebtCeilingTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
