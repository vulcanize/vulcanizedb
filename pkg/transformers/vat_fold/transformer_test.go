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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	vat_fold_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_fold"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
)

type setupOptions struct {
	setMissingHeadersError bool
	setFetcherError        bool
	setConverterError      bool
	setCreateError         bool
	fetchedLogs            []types.Log
	missingHeaders         []core.Header
}

func setup(options setupOptions) (
	vat_fold.VatFoldTransformer,
	*mocks.MockLogFetcher,
	*vat_fold_mocks.MockVatFoldConverter,
	*vat_fold_mocks.MockVatFoldRepository,
) {
	fetcher := &mocks.MockLogFetcher{}
	if options.setFetcherError {
		fetcher.SetFetcherError(fakes.FakeError)
	}
	if len(options.fetchedLogs) > 0 {
		fetcher.SetFetchedLogs(options.fetchedLogs)
	}

	converter := &vat_fold_mocks.MockVatFoldConverter{}
	if options.setConverterError {
		converter.SetConverterError(fakes.FakeError)
	}

	repository := &vat_fold_mocks.MockVatFoldRepository{}
	if options.setMissingHeadersError {
		repository.SetMissingHeadersErr(fakes.FakeError)
	}
	if options.setCreateError {
		repository.SetCreateError(fakes.FakeError)
	}
	if len(options.missingHeaders) > 0 {
		repository.SetMissingHeaders(options.missingHeaders)
	}

	transformer := vat_fold.VatFoldTransformer{
		Config:     vat_fold.VatFoldConfig,
		Fetcher:    fetcher,
		Converter:  converter,
		Repository: repository,
	}

	return transformer, fetcher, converter, repository
}

var _ = Describe("Vat fold transformer", func() {
	It("gets missing headers for block numbers specified in config", func() {
		transformer, _, _, repository := setup(setupOptions{})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_fold.VatFoldConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_fold.VatFoldConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		transformer, _, _, _ := setup(setupOptions{
			setMissingHeadersError: true,
		})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		transformer, fetcher, _, _ := setup(setupOptions{
			missingHeaders: []core.Header{{BlockNumber: 1}, {BlockNumber: 2}},
		})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{vat_fold.VatFoldConfig.ContractAddresses, vat_fold.VatFoldConfig.ContractAddresses}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatFoldSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		transformer, _, _, _ := setup(setupOptions{
			setFetcherError: true,
			missingHeaders:  []core.Header{{BlockNumber: 1}},
		})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		transformer, _, converter, _ := setup(setupOptions{
			fetchedLogs:    []types.Log{test_data.EthVatFoldLog},
			missingHeaders: []core.Header{{BlockNumber: 1}},
		})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthVatFoldLog}))
	})

	It("returns error if converter returns error", func() {
		transformer, _, _, _ := setup(setupOptions{
			setConverterError: true,
			fetchedLogs:       []types.Log{test_data.EthVatFoldLog},
			missingHeaders:    []core.Header{{BlockNumber: 1}},
		})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat fold model", func() {
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		transformer, _, _, repository := setup(setupOptions{
			fetchedLogs:    []types.Log{test_data.EthVatFoldLog},
			missingHeaders: []core.Header{fakeHeader},
		})

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModels).To(Equal([]vat_fold.VatFoldModel{test_data.VatFoldModel}))
	})

	It("returns error if repository returns error for create", func() {
		transformer, _, _, _ := setup(setupOptions{
			fetchedLogs:    []types.Log{test_data.EthVatFoldLog},
			missingHeaders: []core.Header{{BlockNumber: 1, Id: 2}},
			setCreateError: true,
		})

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
