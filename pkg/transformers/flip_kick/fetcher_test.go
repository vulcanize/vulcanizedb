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
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
)

var _ = Describe("Fetcher", func() {
	Describe("FetchLogs", func() {
		var blockChain *fakes.MockBlockChain
		var fetcher flip_kick.Fetcher

		BeforeEach(func() {
			blockChain = fakes.NewMockBlockChain()
			fetcher = flip_kick.Fetcher{Blockchain: blockChain}
		})

		It("fetches logs based on the given query", func() {
			blockNumber := int64(3)
			address := "0x4D2"
			topicZeros := [][]common.Hash{{common.HexToHash("0x")}}

			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(blockNumber),
				ToBlock:   big.NewInt(blockNumber),
				Addresses: []common.Address{common.HexToAddress(address)},
				Topics:    topicZeros,
			}
			fetcher.FetchLogs(address, topicZeros, blockNumber)
			blockChain.AssertGetEthLogsWithCustomQueryCalledWith(query)

		})

		It("returns an error if fetching the logs fails", func() {
			blockChain.SetGetLogsErr(fakes.FakeError)
			_, err := fetcher.FetchLogs("", [][]common.Hash{}, int64(1))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
