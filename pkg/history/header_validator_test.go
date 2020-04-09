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
	"errors"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/pkg/history"
)

var _ = Describe("Header validator", func() {
	var (
		headerRepository *fakes.MockHeaderRepository
		blockChain       *fakes.MockBlockChain
	)

	BeforeEach(func() {
		headerRepository = fakes.NewMockHeaderRepository()
		blockChain = fakes.NewMockBlockChain()
	})

	It("attempts to create every header in the validation window", func() {
		headerRepository.SetMissingBlockNumbers([]int64{})
		blockChain.SetLastBlock(big.NewInt(3))
		validator := history.NewHeaderValidator(blockChain, headerRepository, 2)

		_, err := validator.ValidateHeaders()
		Expect(err).NotTo(HaveOccurred())

		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(3, []int64{1, 2, 3})
	})

	It("propagates header repository errors", func() {
		blockChain.SetLastBlock(big.NewInt(3))
		headerRepositoryError := errors.New("CreateOrUpdate")
		headerRepository.SetCreateOrUpdateHeaderReturnErr(headerRepositoryError)
		validator := history.NewHeaderValidator(blockChain, headerRepository, 2)

		_, err := validator.ValidateHeaders()
		Expect(err.Error()).To(ContainSubstring(headerRepositoryError.Error()))
	})
})
