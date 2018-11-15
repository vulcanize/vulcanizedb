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
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Populating headers", func() {

	var headerRepository *fakes.MockHeaderRepository

	BeforeEach(func() {
		headerRepository = fakes.NewMockHeaderRepository()
	})

	It("returns number of headers added", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(2))
		headerRepository.SetMissingBlockNumbers([]int64{2})

		headersAdded, err := history.PopulateMissingHeaders(blockChain, headerRepository, 1)

		Expect(err).NotTo(HaveOccurred())
		Expect(headersAdded).To(Equal(1))
	})

	It("adds missing headers to the db", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(2))
		headerRepository.SetMissingBlockNumbers([]int64{2})

		_, err := history.PopulateMissingHeaders(blockChain, headerRepository, 1)

		Expect(err).NotTo(HaveOccurred())
		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(1, []int64{2})
	})
})
