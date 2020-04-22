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

package eth_test

import (
	"time"

	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
)

var _ = Describe("StateDiffFetcher", func() {
	Describe("FetchStateDiffsAt", func() {
		var (
			mc               *mocks.BackFillerClient
			stateDiffFetcher *eth.PayloadFetcher
		)
		BeforeEach(func() {
			mc = new(mocks.BackFillerClient)
			setDiffAtErr1 := mc.SetReturnDiffAt(test_data.BlockNumber.Uint64(), test_data.MockStatediffPayload)
			Expect(setDiffAtErr1).ToNot(HaveOccurred())
			setDiffAtErr2 := mc.SetReturnDiffAt(test_data.BlockNumber2.Uint64(), test_data.MockStatediffPayload2)
			Expect(setDiffAtErr2).ToNot(HaveOccurred())
			stateDiffFetcher = eth.NewPayloadFetcher(mc, time.Second*60)
		})
		It("Batch calls statediff_stateDiffAt", func() {
			blockHeights := []uint64{
				test_data.BlockNumber.Uint64(),
				test_data.BlockNumber2.Uint64(),
			}
			stateDiffPayloads, fetchErr := stateDiffFetcher.FetchAt(blockHeights)
			Expect(fetchErr).ToNot(HaveOccurred())
			Expect(len(stateDiffPayloads)).To(Equal(2))
			payload1, ok := stateDiffPayloads[0].(statediff.Payload)
			Expect(ok).To(BeTrue())
			payload2, ok := stateDiffPayloads[1].(statediff.Payload)
			Expect(ok).To(BeTrue())
			Expect(payload1).To(Equal(test_data.MockStatediffPayload))
			Expect(payload2).To(Equal(test_data.MockStatediffPayload2))
		})
	})
})
