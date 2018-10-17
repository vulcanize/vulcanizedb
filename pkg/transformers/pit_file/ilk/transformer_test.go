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

package ilk_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	ilk_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/pit_file/ilk"
)

var _ = Describe("Pit file ilk transformer", func() {
	var (
		config      = ilk.IlkFileConfig
		fetcher     mocks.MockLogFetcher
		converter   ilk_mocks.MockPitFileIlkConverter
		repository  ilk_mocks.MockPitFileIlkRepository
		transformer shared.Transformer
		headerOne   core.Header
		headerTwo   core.Header
	)

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = ilk_mocks.MockPitFileIlkConverter{}
		repository = ilk_mocks.MockPitFileIlkRepository{}
		headerOne = core.Header{Id: GinkgoRandomSeed(), BlockNumber: GinkgoRandomSeed()}
		headerTwo = core.Header{Id: GinkgoRandomSeed(), BlockNumber: GinkgoRandomSeed()}
		transformer = factories.Transformer{
			Config:     config,
			Converter:  &converter,
			Repository: &repository,
			Fetcher:    &fetcher,
		}.NewTransformer(nil, nil)
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
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.PitFileIlkSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{headerOne})

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

	It("converts matching logs", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthPitFileIlkLog}))
	})

	It("returns error if converter returns error", func() {
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists pit file model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository.SetMissingHeaders([]core.Header{headerOne})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels).To(Equal([]interface{}{test_data.PitFileIlkModel}))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileIlkLog})
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetCreateError(fakes.FakeError)

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
