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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ = Describe("FlopKick Converter", func() {
	Describe("ToEntities", func() {
		It("converts a log to a FlopKick entity", func() {
			converter := flop_kick.FlopKickConverter{}
			entities, err := converter.ToEntities(shared.FlopperContractAddress, shared.FlopperABI, []types.Log{test_data.FlopKickLog})

			Expect(err).NotTo(HaveOccurred())
			entity := entities[0]
			Expect(entity.Id).To(Equal(test_data.FlopKickEntity.Id))
			Expect(entity.Lot).To(Equal(test_data.FlopKickEntity.Lot))
			Expect(entity.Bid).To(Equal(test_data.FlopKickEntity.Bid))
			Expect(entity.Gal).To(Equal(test_data.FlopKickEntity.Gal))
			Expect(entity.End).To(Equal(test_data.FlopKickEntity.End))
			Expect(entity.TransactionIndex).To(Equal(test_data.FlopKickEntity.TransactionIndex))
			Expect(entity.Raw).To(Equal(test_data.FlopKickEntity.Raw))
		})

		It("returns an error if converting the log to an entity fails", func() {
			converter := flop_kick.FlopKickConverter{}
			entities, err := converter.ToEntities(shared.FlopperContractAddress, "error abi", []types.Log{test_data.FlopKickLog})

			Expect(err).To(HaveOccurred())
			Expect(entities).To(BeNil())
		})
	})

	Describe("ToModels", func() {
		var emptyAddressHex = "0x0000000000000000000000000000000000000000"
		var emptyString = ""
		var emptyTime = time.Unix(0, 0)
		var emptyEntities = []flop_kick.Entity{flop_kick.Entity{}}

		It("converts an Entity to a Model", func() {
			converter := flop_kick.FlopKickConverter{}
			models, err := converter.ToModels([]flop_kick.Entity{test_data.FlopKickEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(models[0]).To(Equal(test_data.FlopKickModel))
		})

		It("handles nil values", func() {
			converter := flop_kick.FlopKickConverter{}

			models, err := converter.ToModels(emptyEntities)
			model := models[0]
			Expect(err).NotTo(HaveOccurred())
			Expect(model.BidId).To(Equal(emptyString))
			Expect(model.Lot).To(Equal(emptyString))
			Expect(model.Bid).To(Equal(emptyString))
			Expect(model.Gal).To(Equal(emptyAddressHex))
			Expect(model.End).To(Equal(emptyTime))
		})
	})
})
