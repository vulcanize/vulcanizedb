// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
