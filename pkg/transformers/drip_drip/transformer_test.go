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

package drip_drip_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	drip_drip_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/drip_drip"
)

var _ = Describe("Drip drip transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &drip_drip_mocks.MockDripDripRepository{}
		transformer := drip_drip.DripDripTransformer{
			Config:     drip_drip.DripDripConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &drip_drip_mocks.MockDripDripConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(drip_drip.DripDripConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(drip_drip.DripDripConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &drip_drip_mocks.MockDripDripRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := drip_drip.DripDripTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &drip_drip_mocks.MockDripDripConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &drip_drip_mocks.MockDripDripRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := drip_drip.DripDripTransformer{
			Config:     drip_drip.DripDripConfig,
			Fetcher:    fetcher,
			Converter:  &drip_drip_mocks.MockDripDripConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddress).To(Equal(drip_drip.DripDripConfig.ContractAddress))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.DripDripSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &drip_drip_mocks.MockDripDripRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := drip_drip.DripDripTransformer{
			Fetcher:    fetcher,
			Converter:  &drip_drip_mocks.MockDripDripConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &drip_drip_mocks.MockDripDripConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripDripLog})
		repository := &drip_drip_mocks.MockDripDripRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := drip_drip.DripDripTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLog).To(Equal(test_data.EthDripDripLog))
	})

	It("returns error if converter returns error", func() {
		converter := &drip_drip_mocks.MockDripDripConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripDripLog})
		repository := &drip_drip_mocks.MockDripDripRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := drip_drip.DripDripTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists drip drip model", func() {
		converter := &drip_drip_mocks.MockDripDripConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripDripLog})
		repository := &drip_drip_mocks.MockDripDripRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := drip_drip.DripDripTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModel).To(Equal(test_data.DripDripModel))
	})

	It("returns error if repository returns error for create", func() {
		converter := &drip_drip_mocks.MockDripDripConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthDripDripLog})
		repository := &drip_drip_mocks.MockDripDripRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := drip_drip.DripDripTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
