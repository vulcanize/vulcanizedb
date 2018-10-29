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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
)

var _ = Describe("Frob transformer", func() {
	var (
		repository  mocks.MockRepository
		transformer shared.Transformer
		fetcher     mocks.MockLogFetcher
		converter   mocks.MockConverter
		config      = frob.FrobConfig
	)
	BeforeEach(func() {
		repository = mocks.MockRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockConverter{}
		transformer = factories.Transformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}.NewTransformer(nil, nil)
	})

	It("gets missing headers for block numbers specified in config", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(frob.FrobConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(frob.FrobConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{frob.FrobConfig.ContractAddresses, frob.FrobConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.FrobSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

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
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs to entity", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ContractAbi).To(Equal(frob.FrobConfig.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal([]types.Log{test_data.EthFrobLog}))
	})

	It("returns error if converting to entity returns error", func() {
		converter.ToEntitiesError = fakes.FakeError
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts frob entity to model", func() {
		converter.EntitiesToReturn = []interface{}{test_data.FrobEntity}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.EntitiesToConvert[0]).To(Equal(test_data.FrobEntity))
	})

	It("returns error if converting to model returns error", func() {
		converter.ToModelsError = fakes.FakeError
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists frob model", func() {
		converter.ModelsToReturn = []interface{}{test_data.FrobModel}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels[0]).To(Equal(test_data.FrobModel))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthFrobLog})
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
