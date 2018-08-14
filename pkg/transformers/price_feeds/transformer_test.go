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
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	price_feeds2 "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/price_feeds"
	"math/big"
)

var _ = Describe("Price feed transformer", func() {
	It("gets missing headers for price feeds", func() {
		mockRepository := &price_feeds2.MockPriceFeedRepository{}
		transformer := price_feeds.PriceFeedTransformer{
			Config:     price_feeds.PriceFeedConfig,
			Fetcher:    &price_feeds2.MockPriceFeedFetcher{},
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMissingHeadersCalledwith(price_feeds.PriceFeedConfig.StartingBlockNumber, price_feeds.PriceFeedConfig.EndingBlockNumber)
	})

	It("returns error is missing headers call returns err", func() {
		mockRepository := &price_feeds2.MockPriceFeedRepository{}
		mockRepository.SetMissingHeadersErr(fakes.FakeError)
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    &price_feeds2.MockPriceFeedFetcher{},
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		mockRepository := &price_feeds2.MockPriceFeedRepository{}
		blockNumberOne := int64(1)
		blockNumberTwo := int64(2)
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumberOne}, {BlockNumber: blockNumberTwo}})
		mockFetcher := &price_feeds2.MockPriceFeedFetcher{}
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockFetcher.AssertFetchLogValuesCalledWith([]int64{blockNumberOne, blockNumberTwo})
	})

	It("returns err if fetcher returns err", func() {
		mockRepository := &price_feeds2.MockPriceFeedRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		mockFetcher := &price_feeds2.MockPriceFeedFetcher{}
		mockFetcher.SetReturnErr(fakes.FakeError)
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists model converted from log", func() {
		mockRepository := &price_feeds2.MockPriceFeedRepository{}
		headerID := int64(11111)
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: headerID}})
		mockFetcher := &price_feeds2.MockPriceFeedFetcher{}
		blockNumber := uint64(22222)
		txIndex := uint(33333)
		usdValue := int64(44444)
		etherMultiplier, _ := price_feeds.Ether.Int64()
		rawUsdValue := big.NewInt(0)
		rawUsdValue = rawUsdValue.Mul(big.NewInt(usdValue), big.NewInt(etherMultiplier))
		address := common.BytesToAddress([]byte{1, 2, 3, 4, 5})
		fakeLog := types.Log{
			Address:     address,
			Topics:      nil,
			Data:        rawUsdValue.Bytes(),
			BlockNumber: blockNumber,
			TxHash:      common.Hash{},
			TxIndex:     txIndex,
			BlockHash:   common.Hash{},
			Index:       0,
			Removed:     false,
		}
		mockFetcher.SetReturnLogs([]types.Log{fakeLog})
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		expectedModel := price_feeds.PriceFeedModel{
			BlockNumber:       blockNumber,
			HeaderID:          headerID,
			MedianizerAddress: address.Bytes(),
			UsdValue:          fmt.Sprintf("%d", usdValue),
			TransactionIndex:  txIndex,
		}
		mockRepository.AssertCreateCalledWith(expectedModel)
	})

	It("returns error if creating price feed update returns error", func() {
		mockRepository := &price_feeds2.MockPriceFeedRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		mockRepository.SetCreateErr(fakes.FakeError)
		mockFetcher := &price_feeds2.MockPriceFeedFetcher{}
		mockFetcher.SetReturnLogs([]types.Log{{}})
		transformer := price_feeds.PriceFeedTransformer{
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
