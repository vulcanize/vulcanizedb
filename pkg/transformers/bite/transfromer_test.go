/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package bite_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	bite_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/bite"
)

var _ = Describe("Bite Transformer", func() {
	var repository bite_mocks.MockBiteRepository
	var fetcher mocks.MockLogFetcher
	var converter bite_mocks.MockBiteConverter
	var transformer bite.BiteTransformer
	var blockNumber1 = rand.Int63()
	var blockNumber2 = rand.Int63()

	BeforeEach(func() {
		repository = bite_mocks.MockBiteRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = bite_mocks.MockBiteConverter{}

		transformer = bite.BiteTransformer{
			Repository: &repository,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Config:     bite.BiteConfig,
		}

		transformer.SetConfig(bite.BiteConfig)
	})

	It("gets missing headers for blocks in the configured range", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(bite.BiteConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(bite.BiteConfig.EndingBlockNumber))
	})

	It("returns an error if it fails to get missing headers", func() {
		repository.SetMissingHeadersErr(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("fetches eth logs for each missing header", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1}, {BlockNumber: blockNumber2}})
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.BiteSignature)}}
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{blockNumber1, blockNumber2}))
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{bite.BiteConfig.ContractAddresses, bite.BiteConfig.ContractAddresses}))
	})

	It("returns an error if fetching logs fails", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		fetcher.SetFetcherError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		headerID := int64(123)
		repository.SetMissingHeaders([]core.Header{{Id: headerID}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		repository.SetMarkHeaderCheckedErr(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an eth log to an Entity", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthBiteLog})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ConverterAbi).To(Equal(bite.BiteConfig.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal([]types.Log{test_data.EthBiteLog}))
	})

	It("returns an error if converter fails", func() {
		headerId := int64(1)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1, Id: headerId}})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthBiteLog})
		converter.SetConverterError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the bite record", func() {
		headerId := int64(1)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1, Id: headerId}})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthBiteLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerId))
		Expect(repository.PassedBiteModels).To(Equal([]bite.BiteModel{test_data.BiteModel}))
	})

	It("returns error if persisting bite record fails", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: blockNumber1}})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthBiteLog})
		repository.SetCreateError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
