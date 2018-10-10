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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	mock_vat_heal "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_heal"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
)

var _ = Describe("VatHeal Transformer", func() {
	var mockRepository mock_vat_heal.MockVatHealRepository
	var transformer vat_heal.VatHealTransformer
	var fetcher mocks.MockLogFetcher
	var converter mock_vat_heal.MockVatHealConverter

	BeforeEach(func() {
		mockRepository = mock_vat_heal.MockVatHealRepository{}
		fetcher = mocks.MockLogFetcher{}
		converter = mock_vat_heal.MockVatHealConverter{}
		transformer = vat_heal.VatHealTransformer{
			Repository: &mockRepository,
			Config:     vat_heal.VatHealConfig,
			Fetcher:    &fetcher,
			Converter:  &converter,
		}
	})

	It("gets all of the missing header ids", func() {
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(mockRepository.PassedStartingBlockNumber).To(Equal(vat_heal.VatHealConfig.StartingBlockNumber))
		Expect(mockRepository.PassedEndingBlockNumber).To(Equal(vat_heal.VatHealConfig.EndingBlockNumber))
	})

	It("returns and error if getting the missing headers fails", func() {
		mockRepository.SetMissingHeadersErr(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches vat heal logs for the headers", func() {
		header := core.Header{BlockNumber: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedContractAddresses[0]).To(Equal(vat_heal.VatHealConfig.ContractAddresses))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatHealSignature)}}))
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{header.BlockNumber}))
	})

	It("returns and error if fetching the logs fails", func() {
		header := core.Header{BlockNumber: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetcherError(fakes.FakeError)
		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts the logs to models", func() {
		header := core.Header{BlockNumber: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.VatHealLog}))
	})

	It("returns an error if converting fails", func() {
		header := core.Header{BlockNumber: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		converter.SetConverterError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists the vat heal models", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(mockRepository.PassedModels).To(ContainElement(test_data.VatHealModel))
		Expect(mockRepository.PassedHeaderID).To(Equal(header.Id))
	})

	It("returns an error if persisting the vat heal models fails", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		mockRepository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks the header as checked when there are no logs", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(mockRepository.MarkHeaderCheckedPassedHeaderID).To(Equal(header.Id))
	})

	It("doesn't call MarkCheckedHeader when there are logs", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		fetcher.SetFetchedLogs([]types.Log{test_data.VatHealLog})
		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(mockRepository.MarkHeaderCheckedPassedHeaderID).To(Equal(int64(0)))
	})

	It("returns an error if MarkCheckedHeader fails", func() {
		header := core.Header{Id: GinkgoRandomSeed()}
		mockRepository.SetMissingHeaders([]core.Header{header})
		mockRepository.SetMissingHeadersErr(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
