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
