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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"math/rand"
)

var _ = Describe("DealTransformer", func() {
	var config = deal.DealConfig
	var repository mocks.MockRepository
	var fetcher mocks.MockLogFetcher
	var converter mocks.MockLogNoteConverter
	var transformer shared.Transformer
	var headerOne core.Header
	var headerTwo core.Header

	BeforeEach(func() {
		repository = mocks.MockRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockLogNoteConverter{}
		transformer = factories.LogNoteTransformer{
			Config:     config,
			Converter:  &converter,
			Repository: &repository,
			Fetcher:    &fetcher,
		}.NewLogNoteTransformer(nil, nil)
		headerOne = core.Header{BlockNumber: rand.Int63(), Id: rand.Int63()}
		headerTwo = core.Header{BlockNumber: rand.Int63(), Id: rand.Int63()}
	})

	It("sets the blockchain and database", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers", func() {
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns an error if fetching the missing headers fails", func() {
		repository.SetMissingHeadersError(fakes.FakeError)
		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
	})

	It("fetches logs for each missing header", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{
			config.ContractAddresses, config.ContractAddresses}))
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.DealSignature)}}
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
	})

	It("returns an error if fetching logs fails", func() {
		repository.SetMissingHeaders([]core.Header{{}})
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

	It("converts each eth log to a Model", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})

		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.DealLogNote}))
	})

	It("returns an error if converting the eth log fails", func() {
		repository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})
		converter.SetConverterError(fakes.FakeError)

		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists each model as a Deal record", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})
		converter.SetReturnModels([]interface{}{test_data.DealModel})

		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedModels).To(Equal([]interface{}{test_data.DealModel}))
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
	})

	It("returns an error if persisting deal record fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetCreateError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.DealLogNote})

		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
	})
})
