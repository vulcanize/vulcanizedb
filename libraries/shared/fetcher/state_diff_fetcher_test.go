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
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
)

var _ = Describe("StateDiffFetcher", func() {
	Describe("FetchStateDiffsAt", func() {
		var (
			mc               *mocks.BackFillerClient
			stateDiffFetcher fetcher.StateDiffFetcher
		)
		BeforeEach(func() {
			mc = new(mocks.BackFillerClient)
			setDiffAtErr1 := mc.SetReturnDiffAt(test_data.BlockNumber.Uint64(), test_data.MockStatediffPayload)
			Expect(setDiffAtErr1).ToNot(HaveOccurred())
			setDiffAtErr2 := mc.SetReturnDiffAt(test_data.BlockNumber2.Uint64(), test_data.MockStatediffPayload2)
			Expect(setDiffAtErr2).ToNot(HaveOccurred())
			stateDiffFetcher = fetcher.NewStateDiffFetcher(mc)
		})
		It("Batch calls statediff_stateDiffAt", func() {
			blockHeights := []uint64{
				test_data.BlockNumber.Uint64(),
				test_data.BlockNumber2.Uint64(),
			}
			stateDiffPayloads, fetchErr := stateDiffFetcher.FetchStateDiffsAt(blockHeights)
			Expect(fetchErr).ToNot(HaveOccurred())
			Expect(len(stateDiffPayloads)).To(Equal(2))
			// Can only rlp encode the slice of diffs as part of a struct
			// Rlp encoding allows us to compare content of the slices when the order in the slice may vary
			expectedPayloadsStruct := struct {
				payloads []*statediff.Payload
			}{
				[]*statediff.Payload{
					&test_data.MockStatediffPayload,
					&test_data.MockStatediffPayload2,
				},
			}
			expectedPayloadsBytes, rlpErr1 := rlp.EncodeToBytes(expectedPayloadsStruct)
			Expect(rlpErr1).ToNot(HaveOccurred())
			receivedPayloadsStruct := struct {
				payloads []*statediff.Payload
			}{
				stateDiffPayloads,
			}
			receivedPayloadsBytes, rlpErr2 := rlp.EncodeToBytes(receivedPayloadsStruct)
			Expect(rlpErr2).ToNot(HaveOccurred())
			Expect(bytes.Equal(expectedPayloadsBytes, receivedPayloadsBytes)).To(BeTrue())
		})
	})
})
