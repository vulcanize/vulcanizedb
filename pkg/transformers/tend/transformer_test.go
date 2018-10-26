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
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
)

var _ = Describe("Tend LogNoteTransformer", func() {
	var (
		config      = tend.TendConfig
		converter   mocks.MockLogNoteConverter
		repository  mocks.MockRepository
		fetcher     mocks.MockLogFetcher
		transformer shared.Transformer
		headerOne   core.Header
		headerTwo   core.Header
	)

	BeforeEach(func() {
		converter = mocks.MockLogNoteConverter{}
		repository = mocks.MockRepository{}
		fetcher = mocks.MockLogFetcher{}
		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		headerTwo = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		transformer = factories.LogNoteTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}.NewLogNoteTransformer(nil, nil)
	})

	It("sets the blockchain and database", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers for blocks in the configured range", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(tend.TendConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(tend.TendConfig.EndingBlockNumber))
	})

	It("returns an error if it fails to get missing headers", func() {
		repository.SetMissingHeadersError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("fetches eth logs for each missing header", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.TendFunctionSignature)}}
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{tend.TendConfig.ContractAddresses, tend.TendConfig.ContractAddresses}))
	})

	It("returns an error if fetching logs fails", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
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

	It("converts an eth log to a Model", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.TendLogNote}))
	})

	It("returns an error if converter fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})
		converter.SetConverterError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the tend record", func() {
		converter.SetReturnModels([]interface{}{test_data.TendModel})
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels).To(Equal([]interface{}{test_data.TendModel}))
	})

	It("returns error if persisting tend record fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.TendLogNote})
		repository.SetCreateError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
