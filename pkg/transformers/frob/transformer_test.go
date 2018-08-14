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

package frob_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	frob_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/frob"
)

var _ = Describe("Frob transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &frob_mocks.MockFrobRepository{}
		transformer := frob.FrobTransformer{
			Config:     frob.FrobConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &frob_mocks.MockFrobConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(frob.FrobConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(frob.FrobConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &frob_mocks.MockFrobRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := frob.FrobTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &frob_mocks.MockFrobConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &frob_mocks.MockFrobRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := frob.FrobTransformer{
			Fetcher:    fetcher,
			Converter:  &frob_mocks.MockFrobConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddress).To(Equal(frob.FrobConfig.ContractAddresses))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(frob.FrobEventSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &frob_mocks.MockFrobRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := frob.FrobTransformer{
			Fetcher:    fetcher,
			Converter:  &frob_mocks.MockFrobConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &frob_mocks.MockFrobConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository := &frob_mocks.MockFrobRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := frob.FrobTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedContractAddress).To(Equal(frob.FrobConfig.ContractAddresses))
		Expect(converter.PassedContractABI).To(Equal(frob.FrobConfig.ContractAbi))
		Expect(converter.PassedLog).To(Equal(test_data.EthFrobLog))
		Expect(converter.PassedEntity).To(Equal(test_data.FrobEntity))
	})

	It("returns error if converter returns error", func() {
		converter := &frob_mocks.MockFrobConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository := &frob_mocks.MockFrobRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := frob.FrobTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists frob model", func() {
		converter := &frob_mocks.MockFrobConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository := &frob_mocks.MockFrobRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := frob.FrobTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedTransactionIndex).To(Equal(test_data.EthFrobLog.TxIndex))
		Expect(repository.PassedFrobModel).To(Equal(test_data.FrobModel))
	})

	It("returns error if repository returns error for create", func() {
		converter := &frob_mocks.MockFrobConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository := &frob_mocks.MockFrobRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := frob.FrobTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
