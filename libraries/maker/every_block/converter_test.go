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

package every_block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/maker/every_block"
	"github.com/vulcanize/vulcanizedb/libraries/maker/test_data"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var _ = Describe("FlipKickEntity Converter", func() {
	It("converts an Eth Log to and Entity", func() {
		converter := every_block.FlipKickConverter{}
		entity, err := converter.ToEntity(test_data.TemporaryFlipAddress, every_block.FlipperABI, test_data.EthFlipKickLog)

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
	})

	It("returns an error if converting log to entity fails", func() {
		converter := every_block.FlipKickConverter{}
		_, err := converter.ToEntity(test_data.TemporaryFlipAddress, "error abi", test_data.EthFlipKickLog)

		Expect(err).To(HaveOccurred())
	})

	It("converts and Entity to a Model", func() {
		converter := every_block.FlipKickConverter{}
		model, err := converter.ToModel(test_data.FlipKickEntity)
		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal(test_data.FlipKickModel))
	})

	It("handles nil", func() {
		emptyAddressHex := "0x0000000000000000000000000000000000000000"
		emptyByteArrayHex := "0x0000000000000000000000000000000000000000000000000000000000000000"
		emptyString := ""
		emptyTime := time.Unix(0, 0)
		converter := every_block.FlipKickConverter{}
		emptyEntity := every_block.FlipKickEntity{
			Id:  big.NewInt(1),
			Mom: common.Address{},
			Vat: common.Address{},
			Ilk: [32]byte{},
			Lot: nil,
			Bid: nil,
			Guy: common.Address{},
			Gal: common.Address{},
			End: nil,
			Era: nil,
			Lad: common.Address{},
			Tab: nil,
			Raw: types.Log{},
		}
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
	})

	It("returns an error of the flip kick event id is nil", func() {
		converter := every_block.FlipKickConverter{}
		emptyEntity := every_block.FlipKickEntity{}
		_, err := converter.ToModel(emptyEntity)

		Expect(err).To(HaveOccurred())
	})

})
