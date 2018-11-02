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

package flip_kick_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"math/rand"
)

var _ = Describe("FlipKick Transformer", func() {
	var (
		transformer shared.Transformer
		fetcher     mocks.MockLogFetcher
		converter   mocks.MockConverter
		repository  mocks.MockRepository
		blockNumber int64
		headerId    int64
		headers     []core.Header
		logs        []types.Log
		config      = flip_kick.FlipKickConfig
	)

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockConverter{}
		repository = mocks.MockRepository{}
		transformer = factories.Transformer{
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
			Config:     config,
		}

		blockNumber = rand.Int63()
		headerId = rand.Int63()
		logs = []types.Log{test_data.EthFlipKickLog}
		headers = []core.Header{{
			Id:          headerId,
			BlockNumber: blockNumber,
			Hash:        "0x",
			Raw:         nil,
		}}
	})

	It("fetches logs with the configured contract and topic(s) for each block", func() {
		repository.SetMissingHeaders(headers)
		expectedTopics := [][]common.Hash{{common.HexToHash(shared.FlipKickSignature)}}

		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{flip_kick.FlipKickConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal(expectedTopics))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{blockNumber}))
	})

	It("returns an error if the fetcher fails", func() {
		repository.SetMissingHeaders(headers)
		fetcher.SetFetcherError(fakes.FakeError)

		err := transformer.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed"))
	})

	It("marks header checked if no logs returned", func() {
		repository.SetMissingHeaders([]core.Header{{Id: headerId}})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerId)
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{{Id: headerId}})
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts the logs", func() {
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		fetcher.SetFetchedLogs(logs)

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.ContractAbi).To(Equal(flip_kick.FlipKickConfig.ContractAbi))
		Expect(converter.LogsToConvert).To(Equal(logs))
	})

	It("returns an error if converting the geth log fails", func() {
		repository.SetMissingHeaders(headers)
		fetcher.SetFetchedLogs(logs)
		converter.ToEntitiesError = fakes.FakeError

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed"))
	})

	It("persists a flip_kick record", func() {
		repository.SetMissingHeaders(headers)
		fetcher.SetFetchedLogs(logs)
		converter.ModelsToReturn = []interface{}{test_data.FlipKickModel}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerId))
		Expect(repository.PassedModels[0]).To(Equal(test_data.FlipKickModel))
	})

	It("returns an error if persisting a record fails", func() {
		repository.SetCreateError(fakes.FakeError)
		repository.SetMissingHeaders(headers)
		fetcher.SetFetchedLogs(logs)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed"))
	})

	It("returns an error if fetching missing headers fails", func() {
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("gets missing headers for blocks between the configured block number range", func() {
		repository.SetMissingHeaders(headers)
		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(repository.PassedStartingBlockNumber).To(Equal(flip_kick.FlipKickConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(flip_kick.FlipKickConfig.EndingBlockNumber))
	})
})
