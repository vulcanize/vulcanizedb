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

package vat_fold_test

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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
)

var _ = Describe("Vat fold transformer", func() {
	var (
		config      = vat_fold.VatFoldConfig
		fetcher     mocks.MockLogFetcher
		converter   mocks.MockConverter
		repository  mocks.MockRepository
		transformer shared.Transformer
		headerOne   core.Header
		headerTwo   core.Header
	)

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockConverter{}
		repository = mocks.MockRepository{}
		headerOne = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		headerTwo = core.Header{Id: rand.Int63(), BlockNumber: rand.Int63()}
		transformer = factories.Transformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}.NewTransformer(nil, nil)
	})

	It("sets the blockchain and database", func() {
		Expect(fetcher.SetBcCalled).To(BeTrue())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets missing headers for block numbers specified in config", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{config.ContractAddresses, config.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatFoldSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetcherError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatFoldLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthVatFoldLog}))
	})

	It("returns error if converter returns error", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatFoldLog})
		converter.SetConverterError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat fold model", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatFoldLog})
		converter.SetReturnModels([]interface{}{test_data.VatFoldModel})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels).To(Equal([]interface{}{test_data.VatFoldModel}))
	})

	It("returns error if repository returns error for create", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatFoldLog})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks the header as checked when there are no logs", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(headerOne.Id)
	})

	It("doesn't call MarkHeaderChecked when there are logs", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatFoldLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.MarkHeaderCheckedPassedHeaderIDs).To(BeEmpty())
	})

	It("returns an error if MarkHeaderChecked fails", func() {
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetMarkHeaderCheckedError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
