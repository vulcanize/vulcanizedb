// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	price_feeds_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/price_feeds"
)

var _ = Describe("Price feed transformer", func() {
	It("gets missing headers for price feeds", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		transformer := price_feeds.PriceFeedTransformer{
			Config:     price_feeds.PriceFeedConfig,
			Converter:  mockConverter,
			Fetcher:    &price_feeds_mocks.MockPriceFeedFetcher{},
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMissingHeadersCalledwith(price_feeds.PriceFeedConfig.StartingBlockNumber, price_feeds.PriceFeedConfig.EndingBlockNumber)
	})

	It("returns error is missing headers call returns err", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		mockRepository.SetMissingHeadersErr(fakes.FakeError)
		transformer := price_feeds.PriceFeedTransformer{
			Converter:  mockConverter,
			Fetcher:    &price_feeds_mocks.MockPriceFeedFetcher{},
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		blockNumberOne := int64(1)
		blockNumberTwo := int64(2)
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumberOne}, {BlockNumber: blockNumberTwo}})
		mockFetcher := &price_feeds_mocks.MockPriceFeedFetcher{}
		transformer := price_feeds.PriceFeedTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockFetcher.AssertFetchLogValuesCalledWith([]int64{blockNumberOne, blockNumberTwo})
	})

	It("returns err if fetcher returns err", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		mockFetcher := &price_feeds_mocks.MockPriceFeedFetcher{}
		mockFetcher.SetReturnErr(fakes.FakeError)
		transformer := price_feeds.PriceFeedTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts log to a model", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockFetcher := &price_feeds_mocks.MockPriceFeedFetcher{}
		mockFetcher.SetReturnLogs([]types.Log{test_data.EthPriceFeedLog})
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		headerID := int64(11111)
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: headerID}})
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    mockFetcher,
			Converter:  mockConverter,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(mockConverter.PassedHeaderID).To(Equal(headerID))
		Expect(mockConverter.PassedLog).To(Equal(test_data.EthPriceFeedLog))
	})

	It("returns err if converter returns err", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockConverter.SetConverterErr(fakes.FakeError)
		mockFetcher := &price_feeds_mocks.MockPriceFeedFetcher{}
		mockFetcher.SetReturnLogs([]types.Log{test_data.EthPriceFeedLog})
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		headerID := int64(11111)
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: headerID}})
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    mockFetcher,
			Converter:  mockConverter,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists model converted from log", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		headerID := int64(11111)
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: headerID}})
		mockFetcher := &price_feeds_mocks.MockPriceFeedFetcher{}
		mockFetcher.SetReturnLogs([]types.Log{test_data.EthPriceFeedLog})
		transformer := price_feeds.PriceFeedTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertCreateCalledWith(headerID, test_data.PriceFeedModel)
	})

	It("returns error if creating price feed update returns error", func() {
		mockConverter := &price_feeds_mocks.MockPriceFeedConverter{}
		mockRepository := &price_feeds_mocks.MockPriceFeedRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		mockRepository.SetCreateErr(fakes.FakeError)
		mockFetcher := &price_feeds_mocks.MockPriceFeedFetcher{}
		mockFetcher.SetReturnLogs([]types.Log{{}})
		transformer := price_feeds.PriceFeedTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
