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

package chop_lump_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Cat file chop lump converter", func() {
	var converter chop_lump.CatFileChopLumpConverter

	BeforeEach(func() {
		converter = chop_lump.CatFileChopLumpConverter{}
	})

	Context("chop events", func() {
		It("converts a chop log to a model", func() {
			models, err := converter.ToModels([]types.Log{test_data.EthCatFileChopLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(models).To(Equal([]interface{}{test_data.CatFileChopModel}))
		})
	})

	Context("lump events", func() {
		It("converts a lump log to a model", func() {
			models, err := converter.ToModels([]types.Log{test_data.EthCatFileLumpLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(models).To(Equal([]interface{}{test_data.CatFileLumpModel}))
		})
	})

	It("returns err if log is missing topics", func() {
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModels([]types.Log{badLog})
		Expect(err).To(HaveOccurred())
	})
})
