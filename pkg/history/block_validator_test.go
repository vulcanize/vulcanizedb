// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
	"math/big"
)

var _ = Describe("Blocks validator", func() {

	It("calls create or update for all blocks within the window", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(7))
		blocksRepository := fakes.NewMockBlockRepository()
		validator := history.NewBlockValidator(blockChain, blocksRepository, 2)

		window := validator.ValidateBlocks()

		Expect(window).To(Equal(history.ValidationWindow{LowerBound: 5, UpperBound: 7}))
		blocksRepository.AssertCreateOrUpdateBlockCallCountEquals(3)
	})

	It("returns the number of largest block", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(3))
		maxBlockNumber := blockChain.LastBlock()

		Expect(maxBlockNumber.Int64()).To(Equal(int64(3)))
	})
})
