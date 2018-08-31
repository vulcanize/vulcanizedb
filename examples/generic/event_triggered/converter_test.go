// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var expectedBurnModel = event_triggered.BurnModel{
	TokenName:    "Tusd",
	TokenAddress: constants.TusdContractAddress,
	Burner:       "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
	Tokens:       "1097077688018008265106216665536940668749033598146",
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedBurnEntity = event_triggered.BurnEntity{
	TokenName:    "Tusd",
	TokenAddress: common.HexToAddress(constants.TusdContractAddress),
	Burner:       common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Value:        helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedMintModel = event_triggered.MintModel{
	TokenName:    "Tusd",
	TokenAddress: constants.TusdContractAddress,
	Minter:       constants.TusdContractOwner,
	Mintee:       "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
	Tokens:       "1097077688018008265106216665536940668749033598146",
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedMintEntity = event_triggered.MintEntity{
	TokenName:    "Tusd",
	TokenAddress: common.HexToAddress(constants.TusdContractAddress),
	To:           common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Amount:       helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var burnEvent = core.WatchedEvent{
	LogID:       1,
	Name:        constants.BurnEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.BurnEvent.Signature(),
	Topic1:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic2:      "",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var mintEvent = core.WatchedEvent{
	LogID:       1,
	Name:        constants.MintEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.MintEvent.Signature(),
	Topic1:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic2:      "",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var _ = Describe("Transfer Converter", func() {

	var converter *event_triggered.GenericConverter
	var err error

	BeforeEach(func() {
		converter, err = event_triggered.NewGenericConverter(generic.TusdConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	It("converts a watched burn event into a BurnEntity", func() {

		result, err := converter.ToBurnEntity(burnEvent)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&expectedBurnEntity))
	})

	It("converts a BurnEntity to a BurnModel", func() {

		result, err := converter.ToBurnEntity(burnEvent)
		Expect(err).NotTo(HaveOccurred())

		model := converter.ToBurnModel(result)
		Expect(model).To(Equal(&expectedBurnModel))
	})

})

var _ = Describe("Approval Converter", func() {

	var converter *event_triggered.GenericConverter
	var err error

	BeforeEach(func() {
		converter, err = event_triggered.NewGenericConverter(generic.TusdConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	It("converts a watched mint event into a MintEntity", func() {

		result, err := converter.ToMintEntity(mintEvent)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&expectedMintEntity))
	})

	It("converts a MintEntity to a MintModel", func() {

		result, err := converter.ToMintEntity(mintEvent)
		Expect(err).NotTo(HaveOccurred())

		model := converter.ToMintModel(result)
		Expect(model).To(Equal(&expectedMintModel))
	})

})
