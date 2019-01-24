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

package tend_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Tend TendConverter", func() {
	var converter tend.TendConverter

	BeforeEach(func() {
		converter = tend.TendConverter{}
	})

	Describe("ToModels", func() {
		It("converts an eth log to a db model", func() {
			models, err := converter.ToModels([]types.Log{test_data.TendLogNote})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.TendModel))
		})

		It("returns an error if the log data is empty", func() {
			emptyDataLog := test_data.TendLogNote
			emptyDataLog.Data = []byte{}
			_, err := converter.ToModels([]types.Log{emptyDataLog})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("tend log note data is empty"))
		})

		It("returns an error if the expected amount of topics aren't in the log", func() {
			invalidLog := test_data.TendLogNote
			invalidLog.Topics = []common.Hash{}
			_, err := converter.ToModels([]types.Log{invalidLog})

			Expect(err).To(MatchError("tend log does not contain expected topics"))
		})
	})
})
