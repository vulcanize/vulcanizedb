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

package dai_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered/dai"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers"
)

var expectedTransferModel = event_triggered.TransferModel{
	TokenName:    "Dai",
	TokenAddress: constants.DaiContractAddress,
	To:           "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
	From:         "0x000000000000000000000000000000000000Af21",
	Tokens:       "1097077688018008265106216665536940668749033598146",
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedTransferEntity = dai.TransferEntity{
	TokenName:    "Dai",
	TokenAddress: common.HexToAddress(constants.DaiContractAddress),
	Src:          common.HexToAddress("0x000000000000000000000000000000000000Af21"),
	Dst:          common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Wad:          helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedApprovalModel = event_triggered.ApprovalModel{
	TokenName:    "Dai",
	TokenAddress: constants.DaiContractAddress,
	Owner:        "0x000000000000000000000000000000000000Af21",
	Spender:      "0x09BbBBE21a5975cAc061D82f7b843bCE061BA391",
	Tokens:       "1097077688018008265106216665536940668749033598146",
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var expectedApprovalEntity = dai.ApprovalEntity{
	TokenName:    "Dai",
	TokenAddress: common.HexToAddress(constants.DaiContractAddress),
	Src:          common.HexToAddress("0x000000000000000000000000000000000000Af21"),
	Guy:          common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"),
	Wad:          helpers.BigFromString("1097077688018008265106216665536940668749033598146"),
	Block:        5488076,
	TxHash:       "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
}

var transferEvent = core.WatchedEvent{
	LogID:       1,
	Name:        constants.TransferEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.TransferEvent.Signature(),
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var approvalEvent = core.WatchedEvent{
	LogID:       1,
	Name:        constants.ApprovalEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.DaiContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.ApprovalEvent.Signature(),
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var _ = Describe("Transfer Converter", func() {

	var daiConverter *dai.ERC20Converter
	var err error

	BeforeEach(func() {
		daiConverter, err = dai.NewERC20Converter(generic.DaiConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	It("converts a watched transfer event into a TransferEntity", func() {

		result, err := daiConverter.ToTransferEntity(transferEvent)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&expectedTransferEntity))
	})

	It("converts a TransferEntity to an TransferModel", func() {

		result, err := daiConverter.ToTransferEntity(transferEvent)
		Expect(err).NotTo(HaveOccurred())

		model := daiConverter.ToTransferModel(result)
		Expect(model).To(Equal(&expectedTransferModel))
	})

})

var _ = Describe("Approval Converter", func() {

	var daiConverter *dai.ERC20Converter
	var err error

	BeforeEach(func() {
		daiConverter, err = dai.NewERC20Converter(generic.DaiConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	It("converts a watched approval event into a ApprovalEntity", func() {

		result, err := daiConverter.ToApprovalEntity(approvalEvent)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&expectedApprovalEntity))
	})

	It("converts a ApprovalEntity to an ApprovalModel", func() {

		result, err := daiConverter.ToApprovalEntity(approvalEvent)
		Expect(err).NotTo(HaveOccurred())

		model := daiConverter.ToApprovalModel(result)
		Expect(model).To(Equal(&expectedApprovalModel))
	})

})
