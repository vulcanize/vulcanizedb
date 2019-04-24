// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package fetcher_test

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("LogFetcher", func() {
	Describe("FetchLogs", func() {
		It("fetches logs based on the given query", func() {
			blockChain := fakes.NewMockBlockChain()
			logFetcher := fetcher.NewLogFetcher(blockChain)
			header := fakes.FakeHeader

			addresses := []common.Address{
				common.HexToAddress("0xfakeAddress"),
				common.HexToAddress("0xanotherFakeAddress"),
			}

			topicZeros := []common.Hash{common.BytesToHash([]byte{1, 2, 3, 4, 5})}

			_, err := logFetcher.FetchLogs(addresses, topicZeros, header)

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
			logFetcher := fetcher.NewLogFetcher(blockChain)

			_, err := logFetcher.FetchLogs([]common.Address{}, []common.Hash{}, core.Header{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
