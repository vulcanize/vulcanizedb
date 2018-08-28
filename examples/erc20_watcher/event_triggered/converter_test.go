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
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/big"
)

var expectedTransferModel = event_triggered.TransferModel{
	TokenName:    "Dai",
	TokenAddress: "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359",
	To:           "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
	From:         "0x000000000000000000000000000000000000Af21",
	Tokens:       "1097077688018008265106216665536940668749033598146",
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var transferWad, approvalWad big.Int

var expectedTransferEntity = event_triggered.TransferEntity{
	TokenName:    "Dai",
	TokenAddress: common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	Src:          common.HexToAddress("0x000000000000000000000000000000000000Af21"),
	Dst:          common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Wad:          &transferWad,
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedApprovalModel = event_triggered.ApprovalModel{
	TokenName:    "Dai",
	TokenAddress: "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359",
	Owner:        "0x000000000000000000000000000000000000Af21",
	Spender:      "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
	Tokens:       "1097077688018008265106216665536940668749033598146",
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedApprovalEntity = event_triggered.ApprovalEntity{
	TokenName:    "Dai",
	TokenAddress: common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	Src:          common.HexToAddress("0x000000000000000000000000000000000000Af21"),
	Guy:          common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Wad:          &approvalWad,
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var transferEvent = core.WatchedEvent{
	LogID:       1,
	Name:        "Dai",
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.TransferEventSignature,
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var approvalEvent = core.WatchedEvent{
	LogID:       1,
	Name:        "Dai",
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.ApprovalEventSignature,
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var _ = Describe("Transfer Converter", func() {

	daiConverter := event_triggered.NewERC20Converter(generic.DaiConfig)

	It("converts a watched transfer event into a TransferEntity", func() {

		result, err := daiConverter.ToTransferEntity(transferEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.TokenName).To(Equal(expectedTransferEntity.TokenName))
		Expect(result.TokenAddress).To(Equal(expectedTransferEntity.TokenAddress))
		Expect(result.Dst).To(Equal(expectedTransferEntity.Dst))
		Expect(result.Src).To(Equal(expectedTransferEntity.Src))
		Expect(result.Wad).To(Equal(expectedTransferEntity.Wad))
		Expect(result.Block).To(Equal(expectedTransferEntity.Block))
		Expect(result.TxHash).To(Equal(expectedTransferEntity.TxHash))
	})

	It("converts a TransferEntity to an TransferModel", func() {

		result, err := daiConverter.ToTransferEntity(transferEvent)
		Expect(err).NotTo(HaveOccurred())

		model := daiConverter.ToTransferModel(*result)

		Expect(model.TokenName).To(Equal(expectedTransferModel.TokenName))
		Expect(model.TokenAddress).To(Equal(expectedTransferModel.TokenAddress))
		Expect(model.To).To(Equal(expectedTransferModel.To))
		Expect(model.From).To(Equal(expectedTransferModel.From))
		Expect(model.Tokens).To(Equal(expectedTransferModel.Tokens))
		Expect(model.Block).To(Equal(expectedTransferModel.Block))
		Expect(model.TxHash).To(Equal(expectedTransferModel.TxHash))
	})

})

var _ = Describe("Approval Converter", func() {

	daiConverter := event_triggered.NewERC20Converter(generic.DaiConfig)

	It("converts a watched approval event into a ApprovalEntity", func() {

		result, err := daiConverter.ToApprovalEntity(approvalEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.TokenName).To(Equal(expectedApprovalEntity.TokenName))
		Expect(result.TokenAddress).To(Equal(expectedApprovalEntity.TokenAddress))
		Expect(result.Src).To(Equal(expectedApprovalEntity.Src))
		Expect(result.Guy).To(Equal(expectedApprovalEntity.Guy))
		Expect(result.Wad).To(Equal(expectedApprovalEntity.Wad))
		Expect(result.Block).To(Equal(expectedApprovalEntity.Block))
		Expect(result.TxHash).To(Equal(expectedApprovalEntity.TxHash))
	})

	It("converts a ApprovalEntity to an ApprovalModel", func() {

		result, err := daiConverter.ToApprovalEntity(approvalEvent)
		Expect(err).NotTo(HaveOccurred())

		model := daiConverter.ToApprovalModel(*result)

		Expect(model.TokenName).To(Equal(expectedApprovalModel.TokenName))
		Expect(model.TokenAddress).To(Equal(expectedApprovalModel.TokenAddress))
		Expect(model.Owner).To(Equal(expectedApprovalModel.Owner))
		Expect(model.Spender).To(Equal(expectedApprovalModel.Spender))
		Expect(model.Tokens).To(Equal(expectedApprovalModel.Tokens))
		Expect(model.Block).To(Equal(expectedApprovalModel.Block))
		Expect(model.TxHash).To(Equal(expectedApprovalModel.TxHash))
	})

})

func init() {
	transferWad.SetString("1097077688018008265106216665536940668749033598146", 10)
	approvalWad.SetString("1097077688018008265106216665536940668749033598146", 10)
}
