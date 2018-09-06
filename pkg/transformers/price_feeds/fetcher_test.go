// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var _ = Describe("Price fetcher", func() {
	It("gets log value events from price feed medianizers", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryReturnLogs([]types.Log{{}})
		contractAddresses := []string{"pep-contract-address", "pip-contract-address", "rep-contract-address"}
		fetcher := price_feeds.NewPriceFeedFetcher(mockBlockChain, contractAddresses)
		blockNumber := int64(100)

		_, err := fetcher.FetchLogValues(blockNumber)

		Expect(err).NotTo(HaveOccurred())
		var expectedAddresses []common.Address
		for _, address := range contractAddresses {
			expectedAddresses = append(expectedAddresses, common.HexToAddress(address))
		}
		expectedQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(blockNumber),
			ToBlock:   big.NewInt(blockNumber),
			Addresses: expectedAddresses,
			Topics:    [][]common.Hash{{common.HexToHash(shared.LogValueSignature)}},
		}
		mockBlockChain.AssertGetEthLogsWithCustomQueryCalledWith(expectedQuery)
	})

	It("returns error if getting logs fails", func() {
		mockBlockChain := fakes.NewMockBlockChain()
		mockBlockChain.SetGetEthLogsWithCustomQueryErr(fakes.FakeError)
		fetcher := price_feeds.NewPriceFeedFetcher(mockBlockChain, []string{"contract-address"})

		_, err := fetcher.FetchLogValues(100)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
