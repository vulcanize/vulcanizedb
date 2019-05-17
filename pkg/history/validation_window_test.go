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

package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Validation window", func() {
	It("creates a ValidationWindow equal to (HEAD-windowSize, HEAD)", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(5))

		validationWindow, err := history.MakeValidationWindow(blockChain, 2)

		Expect(err).NotTo(HaveOccurred())
		Expect(validationWindow.LowerBound).To(Equal(int64(3)))
		Expect(validationWindow.UpperBound).To(Equal(int64(5)))
	})

	It("returns the window size", func() {
		window := history.ValidationWindow{LowerBound: 1, UpperBound: 3}

		Expect(window.Size()).To(Equal(2))
	})

	It("generates a range of int64s", func() {
		numberOfBlocksCreated := history.MakeRange(0, 5)
		expected := []int64{0, 1, 2, 3, 4, 5}

		Expect(numberOfBlocksCreated).To(Equal(expected))
	})
})
