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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var _ = Describe("Dent Bid BidFetcher", func() {
	var blockChain *fakes.MockBlockChain
	var fetcher shared.IBidFetcher
	var blockNumber = int64(123)
	var address = "0xfakeAddress"
	var contractAbi = "contractAbi"

	BeforeEach(func() {
		blockChain = fakes.NewMockBlockChain()
		fetcher = shared.NewBidFetcher(blockChain)
	})

	It("fetches a bid record for the given id", func() {
		method := "bids"
		bidId := big.NewInt(1)
		result := &shared.Bid{}
		_, err := fetcher.FetchBid(contractAbi, address, blockNumber, bidId)

		Expect(err).NotTo(HaveOccurred())
		blockChain.AssertFetchContractDataCalledWith(contractAbi, address, method, bidId, result, blockNumber)
	})

	It("returns an error if fetching the bid fails", func() {
		blockChain.SetFetchContractDataErr(fakes.FakeError)
		fetcher := shared.NewBidFetcher(blockChain)
		_, err := fetcher.FetchBid(contractAbi, address, blockNumber, nil)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
