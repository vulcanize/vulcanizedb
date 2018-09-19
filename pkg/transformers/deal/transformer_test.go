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

package deal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	deal_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/deal"
	"math/rand"
)

var _ = Describe("DealTransformer", func() {
	var config = deal.Config
	var dealRepository deal_mocks.MockDealRepository
	var fetcher mocks.MockLogFetcher
	var converter deal_mocks.MockDealConverter
	var transformer deal.DealTransformer

	BeforeEach(func() {
		dealRepository = deal_mocks.MockDealRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = deal_mocks.MockDealConverter{}
		transformer = deal.DealTransformer{
			Repository: &dealRepository,
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
		}
	})

	It("gets missing headers", func() {
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(dealRepository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(dealRepository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns an error if fetching the missing headers fails", func() {
		dealRepository.SetMissingHeadersErr(fakes.FakeError)
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("fetches logs for each missing header", func() {
		header1 := core.Header{BlockNumber: rand.Int63()}
		header2 := core.Header{BlockNumber: rand.Int63()}
		dealRepository.SetMissingHeaders([]core.Header{header1, header2})
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedContractAddress).To(Equal(config.ContractAddress))
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.DealSignature)}}
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{header1.BlockNumber, header2.BlockNumber}))
	})

	It("returns an error if fetching logs fails", func() {
		dealRepository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetcherError(fakes.FakeError)
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		mockConverter := &deal_mocks.MockDealConverter{}
		mockRepository := &deal_mocks.MockDealRepository{}
		headerID := int64(123)
		mockRepository.SetMissingHeaders([]core.Header{{Id: headerID}})
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := deal.DealTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &deal_mocks.MockDealConverter{}
		mockRepository := &deal_mocks.MockDealRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := deal.DealTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts each eth log to a Model", func() {
		dealRepository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(converter.LogsToConvert).To(Equal([]types.Log{test_data.DealLogNote}))
	})

	It("returns an error if converting the eth log fails", func() {
		dealRepository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})
		converter.SetConverterError(fakes.FakeError)
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists each model as a Deal record", func() {
		header1 := core.Header{Id: rand.Int63()}
		header2 := core.Header{Id: rand.Int63()}
		dealRepository.SetMissingHeaders([]core.Header{header1, header2})
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(dealRepository.PassedDealModels).To(Equal([]deal.DealModel{test_data.DealModel, test_data.DealModel}))
		Expect(dealRepository.PassedHeaderIDs).To(Equal([]int64{header1.Id, header2.Id}))
	})

	It("returns an error if persisting deal record fails", func() {
		dealRepository.SetMissingHeaders([]core.Header{{}})
		dealRepository.SetCreateError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
	})
})
