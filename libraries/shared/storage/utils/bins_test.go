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

package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

var _ = Describe("GetBlockHeightBins", func() {
	It("splits a block range up into bins", func() {
		var startingBlock uint64 = 1
		var endingBlock uint64 = 10101
		var batchSize uint64 = 100
		blockRangeBins, err := utils.GetBlockHeightBins(startingBlock, endingBlock, batchSize)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(blockRangeBins)).To(Equal(102))
		Expect(blockRangeBins[101]).To(Equal([]uint64{10101}))

		startingBlock = 101
		endingBlock = 10100
		batchSize = 100
		lastBin := make([]uint64, 0)
		for i := 10001; i <= 10100; i++ {
			lastBin = append(lastBin, uint64(i))
		}
		blockRangeBins, err = utils.GetBlockHeightBins(startingBlock, endingBlock, batchSize)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(blockRangeBins)).To(Equal(100))
		Expect(blockRangeBins[99]).To(Equal(lastBin))
	})

	It("throws an error if the starting block is higher than the ending block", func() {
		var startingBlock uint64 = 10102
		var endingBlock uint64 = 10101
		var batchSize uint64 = 100
		_, err := utils.GetBlockHeightBins(startingBlock, endingBlock, batchSize)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ending block number needs to be greater than starting block number"))
	})

	It("throws an error if the batch size is zero", func() {
		var startingBlock uint64 = 1
		var endingBlock uint64 = 10101
		var batchSize uint64 = 0
		_, err := utils.GetBlockHeightBins(startingBlock, endingBlock, batchSize)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("batchsize needs to be greater than zero"))
	})
})
