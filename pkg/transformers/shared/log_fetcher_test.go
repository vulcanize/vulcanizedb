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

package shared_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var _ = Describe("Fetcher", func() {
	Describe("FetchLogs", func() {

		It("fetches logs based on the given query", func() {
			blockChain := fakes.NewMockBlockChain()
			fetcher := shared.NewFetcher(blockChain)
			blockNumber := int64(123)
			address := "0xfakeAddress"
			topicZeros := [][]common.Hash{{common.BytesToHash([]byte{1, 2, 3, 4, 5})}}

			_, err := fetcher.FetchLogs(address, topicZeros, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			expectedQuery := ethereum.FilterQuery{
				FromBlock: big.NewInt(blockNumber),
				ToBlock:   big.NewInt(blockNumber),
				Addresses: []common.Address{common.HexToAddress(address)},
				Topics:    topicZeros,
			}
			blockChain.AssertGetEthLogsWithCustomQueryCalledWith(expectedQuery)

		})

		It("returns an error if fetching the logs fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetGetEthLogsWithCustomQueryErr(fakes.FakeError)
			fetcher := shared.NewFetcher(blockChain)

			_, err := fetcher.FetchLogs("", [][]common.Hash{}, int64(1))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
