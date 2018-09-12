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

package tend_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	tend_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/tend"
)

var _ = Describe("Tend Transformer", func() {
	var repository tend_mocks.MockTendRepository
	var fetcher mocks.MockLogFetcher
	var converter tend_mocks.MockTendConverter
	var transformer tend.TendTransformer
	var blockNumber1 = rand.Int63()
	var blockNumber2 = rand.Int63()

	BeforeEach(func() {
		repository = tend_mocks.MockTendRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = tend_mocks.MockTendConverter{}

		transformer = tend.TendTransformer{
			Repository: &repository,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Config:     tend.TendConfig,
		}
	})

	It("gets missing headers for blocks in the configured range", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(tend.TendConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(tend.TendConfig.EndingBlockNumber))
	})

	It("returns an error if it fails to get missing headers", func() {
		repository.SetMissingHeadersErr(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("fetches eth logs for each missing header", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1}, {BlockNumber: blockNumber2}})
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.TendFunctionSignature)}}
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{blockNumber1, blockNumber2}))
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedContractAddress).To(Equal(shared.FlipperContractAddress))
	})

	It("returns an error if fetching logs fails", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		fetcher.SetFetcherError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an eth log to an Model", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ConverterContract).To(Equal(tend.TendConfig.ContractAddress))
		Expect(converter.ConverterAbi).To(Equal(tend.TendConfig.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal([]types.Log{test_data.TendLogNote}))
	})

	It("returns an error if converter fails", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})
		converter.SetConverterError(fakes.FakeError)
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the tend record", func() {
		headerId := int64(1)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1, Id: headerId}})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerId))
		Expect(repository.PassedTendModel).To(Equal(test_data.TendModel))
	})

	It("returns error if persisting tend record fails", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1}})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})
		repository.SetCreateError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
