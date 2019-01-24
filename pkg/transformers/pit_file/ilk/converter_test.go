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

package ilk_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Pit file ilk converter", func() {
	It("returns err if log is missing topics", func() {
		converter := ilk.PitFileIlkConverter{}
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		converter := ilk.PitFileIlkConverter{}
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("returns error if 'what' field is unknown", func() {
		log := types.Log{
			Address: test_data.EthPitFileIlkLineLog.Address,
			Topics: []common.Hash{
				test_data.EthPitFileIlkLineLog.Topics[0],
				test_data.EthPitFileIlkLineLog.Topics[1],
				test_data.EthPitFileIlkLineLog.Topics[2],
				common.HexToHash("0x1111111100000000000000000000000000000000000000000000000000000000"),
			},
			Data:        test_data.EthPitFileIlkLineLog.Data,
			BlockNumber: test_data.EthPitFileIlkLineLog.BlockNumber,
			TxHash:      test_data.EthPitFileIlkLineLog.TxHash,
			TxIndex:     test_data.EthPitFileIlkLineLog.TxIndex,
			BlockHash:   test_data.EthPitFileIlkLineLog.BlockHash,
			Index:       test_data.EthPitFileIlkLineLog.Index,
		}
		converter := ilk.PitFileIlkConverter{}

		_, err := converter.ToModels([]types.Log{log})

		Expect(err).To(HaveOccurred())
	})

	Describe("when log is valid", func() {
		It("converts to model with data converted to ray when what is 'spot'", func() {
			converter := ilk.PitFileIlkConverter{}

			models, err := converter.ToModels([]types.Log{test_data.EthPitFileIlkSpotLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0].(ilk.PitFileIlkModel)).To(Equal(test_data.PitFileIlkSpotModel))
		})

		It("converts to model with data converted to wad when what is 'line'", func() {
			converter := ilk.PitFileIlkConverter{}

			models, err := converter.ToModels([]types.Log{test_data.EthPitFileIlkLineLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0].(ilk.PitFileIlkModel)).To(Equal(test_data.PitFileIlkLineModel))
		})
	})

})
