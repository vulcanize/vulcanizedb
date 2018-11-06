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

package factories_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"math/rand"
)

var _ = Describe("Transformer", func() {
	var (
		repository  mocks.MockRepository
		fetcher     mocks.MockLogFetcher
		converter   mocks.MockConverter
		transformer shared.Transformer
		headerOne   core.Header
		headerTwo   core.Header
		config      = test_data.GenericTestConfig
		logs        = test_data.GenericTestLogs
	)

	BeforeEach(func() {
		repository = mocks.MockRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockConverter{}

		transformer = factories.Transformer{
			Repository: &repository,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Config:     config,
		}.NewTransformer(nil, nil)

		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		headerTwo = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
	})

	It("sets the blockchain and db", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers for blocks in the configured range", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns an error if it fails to get missing headers", func() {
		repository.SetMissingHeadersError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches eth logs for each missing header", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		expectedTopics := [][]common.Hash{{common.HexToHash(config.Topic)}}
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{config.ContractAddresses, config.ContractAddresses}))
	})

	It("returns an error if fetching logs fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
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

	It("doesn't attempt to convert or persist an empty collection when there are no logs", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ToEntitiesCalledCounter).To(Equal(0))
		Expect(converter.ToModelsCalledCounter).To(Equal(0))
		Expect(repository.CreateCalledCounter).To(Equal(0))
	})

	It("does not call repository.MarkCheckedHeader when there are logs", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedNotCalled()
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an eth log to an entity", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ContractAbi).To(Equal(config.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal(logs))
	})

	It("returns an error if converter fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)
		converter.ToEntitiesError = fakes.FakeError

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts an entity to a model", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)
		converter.EntitiesToReturn = []interface{}{test_data.GenericEntity{}}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.EntitiesToConvert[0]).To(Equal(test_data.GenericEntity{}))
	})

	It("returns an error if converting to models fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)
		converter.EntitiesToReturn = []interface{}{test_data.GenericEntity{}}
		converter.ToModelsError = fakes.FakeError

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the record", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)
		converter.ModelsToReturn = []interface{}{test_data.GenericModel{}}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels[0]).To(Equal(test_data.GenericModel{}))
	})

	It("returns error if persisting the record fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs(logs)
		repository.SetCreateError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
