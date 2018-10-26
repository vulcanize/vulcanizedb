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

package vat_heal_test

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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
)

var _ = Describe("VatHeal Transformer", func() {
	var (
		repository  mocks.MockRepository
		transformer shared.Transformer
		fetcher     mocks.MockLogFetcher
		converter   mocks.MockConverter
		config      = vat_heal.VatHealConfig
	)

	BeforeEach(func() {
		repository = mocks.MockRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = mocks.MockConverter{}
		transformer = factories.Transformer{
			Repository: &repository,
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
		}.NewTransformer(nil, nil)
	})

	It("sets the database and blockchain", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.SetDbCalled).To(BeTrue())
	})

	It("gets all of the missing header ids", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(config.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(config.EndingBlockNumber))
	})

	It("returns and error if getting the missing headers fails", func() {
		repository.SetMissingHeadersError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches vat heal logs for the headers", func() {
		header := core.Header{BlockNumber: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedContractAddresses[0]).To(Equal(vat_heal.VatHealConfig.ContractAddresses))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatHealSignature)}}))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{header.BlockNumber}))
	})

	It("returns and error if fetching the logs fails", func() {
		header := core.Header{BlockNumber: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetcherError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts the logs to models", func() {
		header := core.Header{BlockNumber: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.VatHealLog}))
	})

	It("returns an error if converting fails", func() {
		header := core.Header{BlockNumber: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		converter.SetConverterError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the vat heal models", func() {
		header := core.Header{Id: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		converter.SetReturnModels([]interface{}{test_data.VatHealModel})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedModels).To(ContainElement(test_data.VatHealModel))
		Expect(repository.PassedHeaderID).To(Equal(header.Id))
	})

	It("returns an error if persisting the vat heal models fails", func() {
		header := core.Header{Id: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks the header as checked when there are no logs", func() {
		header := core.Header{Id: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedCalledWith(header.Id)
	})

	It("doesn't call MarkCheckedHeader when there are logs", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		repository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		repository.AssertMarkHeaderCheckedNotCalled()
	})

	It("returns an error if MarkCheckedHeader fails", func() {
		header := core.Header{Id: rand.Int63()}
		repository.SetMissingHeaders([]core.Header{header})
		repository.SetMissingHeadersError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
