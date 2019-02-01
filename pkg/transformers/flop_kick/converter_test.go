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

package flop_kick_test

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("FlopKick Converter", func() {
	Describe("ToEntities", func() {
		It("converts a log to a FlopKick entity", func() {
			converter := flop_kick.FlopKickConverter{}
			entities, err := converter.ToEntities(test_data.KovanFlopperABI, []types.Log{test_data.FlopKickLog})

			Expect(err).NotTo(HaveOccurred())
			entity := entities[0]
			Expect(entity).To(Equal(test_data.FlopKickEntity))
		})

		It("returns an error if converting the log to an entity fails", func() {
			converter := flop_kick.FlopKickConverter{}
			entities, err := converter.ToEntities("error abi", []types.Log{test_data.FlopKickLog})

			Expect(err).To(HaveOccurred())
			Expect(entities).To(BeNil())
		})
	})

	Describe("ToModels", func() {
		var emptyAddressHex = "0x0000000000000000000000000000000000000000"
		var emptyString = ""
		var emptyTime = time.Unix(0, 0)
		var emptyEntity = flop_kick.Entity{}

		It("converts an Entity to a Model", func() {
			converter := flop_kick.FlopKickConverter{}
			models, err := converter.ToModels([]interface{}{test_data.FlopKickEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(models[0]).To(Equal(test_data.FlopKickModel))
		})

		It("returns error if wrong entity", func() {
			converter := flop_kick.FlopKickConverter{}
			_, err := converter.ToModels([]interface{}{test_data.WrongEntity{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("entity of type test_data.WrongEntity, not flop_kick.Entity"))
		})

		It("handles nil values", func() {
			emptyLog, err := json.Marshal(types.Log{})

			converter := flop_kick.FlopKickConverter{}
			expectedModel := flop_kick.Model{
				BidId:            emptyString,
				Lot:              emptyString,
				Bid:              emptyString,
				Gal:              emptyAddressHex,
				End:              emptyTime,
				TransactionIndex: 0,
				Raw:              emptyLog,
			}

			models, err := converter.ToModels([]interface{}{emptyEntity})
			model := models[0]
			Expect(err).NotTo(HaveOccurred())
			Expect(model).To(Equal(expectedModel))
		})
	})
})
