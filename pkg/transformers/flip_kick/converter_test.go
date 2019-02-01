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

package flip_kick_test

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("FlipKick Converter", func() {
	var converter = flip_kick.FlipKickConverter{}

	Describe("ToEntity", func() {
		It("converts an Eth Log to a FlipKickEntity", func() {
			entities, err := converter.ToEntities(test_data.KovanFlipperABI, []types.Log{test_data.EthFlipKickLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			entity := entities[0]
			Expect(entity).To(Equal(test_data.FlipKickEntity))
		})

		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToEntities("error abi", []types.Log{test_data.EthFlipKickLog})

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ToModel", func() {
		var emptyAddressHex = "0x0000000000000000000000000000000000000000"
		var emptyString = ""
		var emptyTime = time.Unix(0, 0)
		var emptyEntity = flip_kick.FlipKickEntity{}
		var emptyRawLog []byte
		var err error

		BeforeEach(func() {
			emptyEntity.Id = big.NewInt(1)
			emptyRawLog, err = json.Marshal(types.Log{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("converts an Entity to a Model", func() {
			models, err := converter.ToModels([]interface{}{test_data.FlipKickEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.FlipKickModel))
		})

		It("handles nil values", func() {
			models, err := converter.ToModels([]interface{}{emptyEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			model := models[0].(flip_kick.FlipKickModel)
			Expect(model.BidId).To(Equal("1"))
			Expect(model.Lot).To(Equal(emptyString))
			Expect(model.Bid).To(Equal(emptyString))
			Expect(model.Gal).To(Equal(emptyAddressHex))
			Expect(model.End).To(Equal(emptyTime))
			Expect(model.Urn).To(Equal(emptyAddressHex))
			Expect(model.Tab).To(Equal(emptyString))
			Expect(model.Raw).To(Equal(emptyRawLog))
		})

		It("returns an error if the flip kick event id is nil", func() {
			emptyEntity.Id = nil
			_, err := converter.ToModels([]interface{}{emptyEntity})

			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the wrong entity type is passed in", func() {
			_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("entity of type"))
		})
	})
})
