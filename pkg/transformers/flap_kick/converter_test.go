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

package flap_kick_test

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flap_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"math/big"
	"time"
)

var _ = Describe("Flap kick converter", func() {
	var converter = flap_kick.FlapKickConverter{}

	Describe("ToEntity", func() {
		It("converts an Eth Log to a FlapKickEntity", func() {
			entities, err := converter.ToEntities(test_data.KovanFlapperABI, []types.Log{test_data.EthFlapKickLog})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0]).To(Equal(test_data.FlapKickEntity))
		})

		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToEntities("error abi", []types.Log{test_data.EthFlapKickLog})

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ToModel", func() {

		BeforeEach(func() {
		})

		It("converts an Entity to a Model", func() {
			models, err := converter.ToModels([]interface{}{test_data.FlapKickEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.FlapKickModel))
		})

		It("handles nil values", func() {
			emptyAddressHex := "0x0000000000000000000000000000000000000000"
			emptyString := ""
			emptyTime := time.Unix(0, 0)
			emptyEntity := flap_kick.FlapKickEntity{}
			emptyEntity.Id = big.NewInt(1)
			emptyRawLogJson, err := json.Marshal(types.Log{})
			Expect(err).NotTo(HaveOccurred())

			models, err := converter.ToModels([]interface{}{emptyEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			model := models[0].(flap_kick.FlapKickModel)
			Expect(model.BidId).To(Equal("1"))
			Expect(model.Lot).To(Equal(emptyString))
			Expect(model.Bid).To(Equal(emptyString))
			Expect(model.Gal).To(Equal(emptyAddressHex))
			Expect(model.End).To(Equal(emptyTime))
			Expect(model.Raw).To(Equal(emptyRawLogJson))
		})

		It("returns an error if the flap kick event id is nil", func() {
			_, err := converter.ToModels([]interface{}{flap_kick.FlapKickEntity{}})

			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the wrong entity type is passed in", func() {
			_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("entity of type"))
		})
	})
})
