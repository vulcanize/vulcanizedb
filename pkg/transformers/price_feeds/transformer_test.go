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

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"math/rand"
)

var _ = Describe("Price feed transformer", func() {
	var (
		config      = price_feeds.PriceFeedConfig
		fetcher     mocks.MockLogFetcher
		converter   mocks.MockLogNoteConverter
		repository  mocks.MockRepository
		transformer shared.Transformer
		headerOne   core.Header
		headerTwo   core.Header
	)

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockLogNoteConverter{}
		repository = mocks.MockRepository{}
		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		headerTwo = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		transformer = factories.LogNoteTransformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}.NewLogNoteTransformer(nil, nil)
	})

	It("sets the blockchain and db", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers for price feeds", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns error is missing headers call returns err", func() {
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{config.ContractAddresses, config.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(config.Topic)}}))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
	})

	It("returns err if fetcher returns err", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetcherError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerOne.Id)
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts log to a model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPriceFeedLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthPriceFeedLog}))
	})

	It("returns err if converter returns err", func() {
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPriceFeedLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists model converted from log", func() {
		converter.SetReturnModels([]interface{}{test_data.PriceFeedModel})
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPriceFeedLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(len(repository.PassedModels)).To(Equal(1))
		Expect(repository.PassedModels[0]).To(Equal(test_data.PriceFeedModel))
	})

	It("returns error if creating price feed update returns error", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetCreateError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{{}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
