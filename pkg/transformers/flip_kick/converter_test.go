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

package flip_kick_test

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("FlipKick Converter", func() {
	var converter = flip_kick.FlipKickConverter{}

	Describe("ToEntity", func() {
		It("converts an Eth Log to a FlipKickEntity", func() {
			entities, err := converter.ToEntities(shared.FlipperABI, []types.Log{test_data.EthFlipKickLog})

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
		var emptyRawLog string

		BeforeEach(func() {
			emptyEntity.Id = big.NewInt(1)
			var emptyRawLogJson, err = json.Marshal(types.Log{})
			Expect(err).NotTo(HaveOccurred())

			emptyRawLogJson, err = json.Marshal(types.Log{})
			Expect(err).NotTo(HaveOccurred())
			emptyRawLog = string(emptyRawLogJson)
		})

		It("converts an Entity to a Model", func() {
			models, err := converter.ToModels([]flip_kick.FlipKickEntity{test_data.FlipKickEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			Expect(models[0]).To(Equal(test_data.FlipKickModel))
		})

		It("handles nil values", func() {
			models, err := converter.ToModels([]flip_kick.FlipKickEntity{emptyEntity})

			Expect(err).NotTo(HaveOccurred())
			Expect(len(models)).To(Equal(1))
			model := models[0]
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
			_, err := converter.ToModels([]flip_kick.FlipKickEntity{emptyEntity})

			Expect(err).To(HaveOccurred())
		})
	})
})
