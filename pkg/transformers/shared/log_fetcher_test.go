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
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var _ = Describe("Fetcher", func() {
	Describe("FetchLogs", func() {
		It("fetches logs based on the given query", func() {
			blockChain := fakes.NewMockBlockChain()
			fetcher := shared.NewFetcher(blockChain)
			header := fakes.FakeHeader

			addresses := []common.Address{
				common.HexToAddress("0xfakeAddress"),
				common.HexToAddress("0xanotherFakeAddress"),
			}

			topicZeros := []common.Hash{common.BytesToHash([]byte{1, 2, 3, 4, 5})}

			_, err := fetcher.FetchLogs(addresses, topicZeros, header)

			address1 := common.HexToAddress("0xfakeAddress")
			address2 := common.HexToAddress("0xanotherFakeAddress")
			Expect(err).NotTo(HaveOccurred())

			blockHash := common.HexToHash(header.Hash)
			expectedQuery := ethereum.FilterQuery{
				BlockHash: &blockHash,
				Addresses: []common.Address{address1, address2},
				Topics:    [][]common.Hash{topicZeros},
			}
			blockChain.AssertGetEthLogsWithCustomQueryCalledWith(expectedQuery)
		})

		It("returns an error if fetching the logs fails", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetGetEthLogsWithCustomQueryErr(fakes.FakeError)
			fetcher := shared.NewFetcher(blockChain)

			_, err := fetcher.FetchLogs([]common.Address{}, []common.Hash{}, core.Header{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
