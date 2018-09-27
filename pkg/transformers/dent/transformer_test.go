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

package dent_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	dent_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/dent"
)

var _ = Describe("DentTransformer", func() {
	var config = dent.DentConfig
	var dentRepository dent_mocks.MockDentRepository
	var fetcher mocks.MockLogFetcher
	var converter dent_mocks.MockDentConverter
	var transformer dent.DentTransformer

	BeforeEach(func() {
		dentRepository = dent_mocks.MockDentRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = dent_mocks.MockDentConverter{}
		transformer = dent.DentTransformer{
			Repository: &dentRepository,
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
		}
	})

	It("gets missing headers", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(dentRepository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(dentRepository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns an error if fetching the missing headers fails", func() {
		dentRepository.SetMissingHeadersError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for each missing header", func() {
		header1 := core.Header{BlockNumber: rand.Int63()}
		header2 := core.Header{BlockNumber: rand.Int63()}
		dentRepository.SetMissingHeaders([]core.Header{header1, header2})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{config.ContractAddresses, config.ContractAddresses}))
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.DentFunctionSignature)}}
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{header1.BlockNumber, header2.BlockNumber}))
	})

	It("returns an error if fetching logs fails", func() {
		dentRepository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetcherError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header checked if no logs returned", func() {
		mockConverter := &dent_mocks.MockDentConverter{}
		mockRepository := &dent_mocks.MockDentRepository{}
		headerID := int64(123)
		mockRepository.SetMissingHeaders([]core.Header{{Id: headerID}})
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := dent.DentTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		mockRepository.AssertMarkHeaderCheckedCalledWith(headerID)
	})

	It("returns error if marking header checked returns err", func() {
		mockConverter := &dent_mocks.MockDentConverter{}
		mockRepository := &dent_mocks.MockDentRepository{}
		mockRepository.SetMissingHeaders([]core.Header{{Id: int64(123)}})
		mockRepository.SetMarkHeaderCheckedErr(fakes.FakeError)
		mockFetcher := &mocks.MockLogFetcher{}
		transformer := dent.DentTransformer{
			Converter:  mockConverter,
			Fetcher:    mockFetcher,
			Repository: mockRepository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts each eth log to a Model", func() {
		dentRepository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.DentLog})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.LogsToConvert).To(Equal([]types.Log{test_data.DentLog}))
	})

	It("returns an error if converting the eth log fails", func() {
		dentRepository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.DentLog})
		converter.SetConverterError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists each model as a Dent record", func() {
		header1 := core.Header{Id: rand.Int63()}
		header2 := core.Header{Id: rand.Int63()}
		dentRepository.SetMissingHeaders([]core.Header{header1, header2})
		fetcher.SetFetchedLogs([]types.Log{test_data.DentLog})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(dentRepository.PassedDentModels).To(Equal([]dent.DentModel{test_data.DentModel, test_data.DentModel}))
		Expect(dentRepository.PassedHeaderIds).To(Equal([]int64{header1.Id, header2.Id}))
	})

	It("returns an error if persisting dent record fails", func() {
		dentRepository.SetMissingHeaders([]core.Header{{}})
		dentRepository.SetCreateError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.DentLog})
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
	})
})
