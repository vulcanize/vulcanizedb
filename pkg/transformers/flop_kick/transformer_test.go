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

package flop_kick_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	flop_kick_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/flop_kick"
)

var _ = Describe("FlopKick Transformer", func() {
	var fetcher mocks.MockLogFetcher
	var converter flop_kick_mocks.MockConverter
	var repository flop_kick_mocks.MockRepository
	var config = flop_kick.Config
	var headerOne core.Header
	var headerTwo core.Header

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = flop_kick_mocks.MockConverter{}
		repository = flop_kick_mocks.MockRepository{}
		headerOne = core.Header{Id: GinkgoRandomSeed(), BlockNumber: GinkgoRandomSeed()}
		headerTwo = core.Header{Id: GinkgoRandomSeed(), BlockNumber: GinkgoRandomSeed()}
	})

	It("gets missing headers for specified block numbers", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}

		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(flop_kick.Config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(flop_kick.Config.EndingBlockNumber))
	})

	It("returns error if getting missing headers fails", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for each missing header", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedContractAddress).To(Equal(flop_kick.Config.ContractAddress))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.FlopKickSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{{}})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header as checked even if no logs were returned", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		fetcher.SetFetchedLogs([]types.Log{})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.CheckedHeaderIds).To(ContainElement(headerOne.Id))
		Expect(repository.CheckedHeaderIds).To(ContainElement(headerTwo.Id))
	})

	It("returns error if marking header checked returns err", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		repository.SetCheckedHeaderError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs to entity", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: GinkgoRandomSeed()}})
		fetcher.SetFetchedLogs([]types.Log{test_data.FlopKickLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedContractAddress).To(Equal(flop_kick.Config.ContractAddress))
		Expect(converter.PassedContractABI).To(Equal(flop_kick.Config.ContractAbi))
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.FlopKickLog}))
	})

	It("returns an error if converting logs to entity fails", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: GinkgoRandomSeed()}})
		fetcher.SetFetchedLogs([]types.Log{test_data.FlopKickLog})
		converter.SetToEntityConverterError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts flop_kick entity to model", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.FlopKickLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedEntities).To(Equal([]flop_kick.Entity{test_data.FlopKickEntity}))
	})

	It("returns an error if there's a failure in converting to model", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{{}})
		fetcher.SetFetchedLogs([]types.Log{test_data.FlopKickLog})
		converter.SetToModelConverterError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the flop_kick model", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		fetcher.SetFetchedLogs([]types.Log{test_data.FlopKickLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.CreatedHeaderIds).To(ContainElement(headerOne.Id))
		Expect(repository.CreatedHeaderIds).To(ContainElement(headerTwo.Id))
		Expect(repository.CreatedModels).To(ContainElement(test_data.FlopKickModel))
	})

	It("returns error if repository returns error for create", func() {
		transformer := flop_kick.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}
		repository.SetMissingHeaders([]core.Header{{}})
		repository.SetCreateError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.FlopKickLog})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
