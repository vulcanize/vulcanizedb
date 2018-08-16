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
			entity, err := converter.ToEntity(test_data.FlipAddress, shared.FlipperABI, test_data.EthFlipKickLog)

			Expect(err).NotTo(HaveOccurred())
			Expect(entity.Id).To(Equal(test_data.FlipKickEntity.Id))
			Expect(entity.Mom).To(Equal(test_data.FlipKickEntity.Mom))
			Expect(entity.Vat).To(Equal(test_data.FlipKickEntity.Vat))
			Expect(entity.Ilk).To(Equal(test_data.FlipKickEntity.Ilk))
			Expect(entity.Lot).To(Equal(test_data.FlipKickEntity.Lot))
			Expect(entity.Bid).To(Equal(test_data.FlipKickEntity.Bid))
			Expect(entity.Guy).To(Equal(test_data.FlipKickEntity.Guy))
			Expect(entity.Gal).To(Equal(test_data.FlipKickEntity.Gal))
			Expect(entity.End).To(Equal(test_data.FlipKickEntity.End))
			Expect(entity.Era).To(Equal(test_data.FlipKickEntity.Era))
			Expect(entity.Lad).To(Equal(test_data.FlipKickEntity.Lad))
			Expect(entity.Tab).To(Equal(test_data.FlipKickEntity.Tab))
			Expect(entity.Raw).To(Equal(test_data.FlipKickEntity.Raw))
		})

		It("returns an error if converting log to entity fails", func() {
			_, err := converter.ToEntity(test_data.FlipAddress, "error abi", test_data.EthFlipKickLog)

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ToModel", func() {
		var emptyAddressHex = "0x0000000000000000000000000000000000000000"
		var emptyByteArrayHex = "0x0000000000000000000000000000000000000000000000000000000000000000"
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
			model, err := converter.ToModel(test_data.FlipKickEntity)

			Expect(err).NotTo(HaveOccurred())
			Expect(model).To(Equal(test_data.FlipKickModel))
		})

		It("handles nil values", func() {
			model, err := converter.ToModel(emptyEntity)

			Expect(err).NotTo(HaveOccurred())
			Expect(model.Id).To(Equal("1"))
			Expect(model.Mom).To(Equal(emptyAddressHex))
			Expect(model.Vat).To(Equal(emptyAddressHex))
			Expect(model.Ilk).To(Equal(emptyByteArrayHex))
			Expect(model.Lot).To(Equal(emptyString))
			Expect(model.Bid).To(Equal(emptyString))
			Expect(model.Guy).To(Equal(emptyAddressHex))
			Expect(model.Gal).To(Equal(emptyAddressHex))
			Expect(model.End).To(Equal(emptyTime))
			Expect(model.Era).To(Equal(emptyTime))
			Expect(model.Lad).To(Equal(emptyAddressHex))
			Expect(model.Tab).To(Equal(emptyString))
			Expect(model.Raw).To(Equal(emptyRawLog))
		})

		It("returns an error if the flip kick event id is nil", func() {
			emptyEntity.Id = nil
			entity, err := converter.ToModel(emptyEntity)

			Expect(err).To(HaveOccurred())
			Expect(entity).To(Equal(flip_kick.FlipKickModel{}))
		})
	})
})
