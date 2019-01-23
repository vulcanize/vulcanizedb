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

package common_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

var _ = Describe("Conversion of GethReceipt to core.Receipt", func() {

	It(`converts geth receipt to internal receipt format (pre Byzantium has post-transaction stateroot)`, func() {
		receipt := types.Receipt{
			Bloom:             types.BytesToBloom(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			ContractAddress:   common.Address{},
			CumulativeGasUsed: uint64(25000),
			GasUsed:           uint64(21000),
			Logs:              []*types.Log{},
			PostState:         hexutil.MustDecode("0x88abf7e73128227370aa7baa3dd4e18d0af70e92ef1f9ef426942fbe2dddb733"),
			TxHash:            common.HexToHash("0x97d99bc7729211111a21b12c933c949d4f31684f1d6954ff477d0477538ff017"),
		}

		expected := core.Receipt{
			Bloom:             "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ContractAddress:   "",
			CumulativeGasUsed: 25000,
			GasUsed:           21000,
			Logs:              []core.Log{},
			StateRoot:         "0x88abf7e73128227370aa7baa3dd4e18d0af70e92ef1f9ef426942fbe2dddb733",
			Status:            -99,
			TxHash:            receipt.TxHash.Hex(),
		}

		coreReceipt := vulcCommon.ToCoreReceipt(&receipt)
		Expect(coreReceipt.Bloom).To(Equal(expected.Bloom))
		Expect(coreReceipt.ContractAddress).To(Equal(expected.ContractAddress))
		Expect(coreReceipt.CumulativeGasUsed).To(Equal(expected.CumulativeGasUsed))
		Expect(coreReceipt.GasUsed).To(Equal(expected.GasUsed))
		Expect(coreReceipt.Logs).To(Equal(expected.Logs))
		Expect(coreReceipt.StateRoot).To(Equal(expected.StateRoot))
		Expect(coreReceipt.Status).To(Equal(expected.Status))
		Expect(coreReceipt.TxHash).To(Equal(expected.TxHash))

	})

	It("converts geth receipt to internal receipt format (post Byzantium has status", func() {
		receipt := types.Receipt{
			Bloom:             types.BytesToBloom(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			ContractAddress:   common.HexToAddress("x0123"),
			CumulativeGasUsed: uint64(7996119),
			GasUsed:           uint64(21000),
			Logs:              []*types.Log{},
			Status:            uint64(1),
			TxHash:            common.HexToHash("0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223"),
		}

		expected := core.Receipt{
			Bloom:             "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ContractAddress:   receipt.ContractAddress.Hex(),
			CumulativeGasUsed: 7996119,
			GasUsed:           21000,
			Logs:              []core.Log{},
			StateRoot:         "",
			Status:            1,
			TxHash:            receipt.TxHash.Hex(),
		}

		coreReceipt := vulcCommon.ToCoreReceipt(&receipt)
		Expect(coreReceipt.Bloom).To(Equal(expected.Bloom))
		Expect(coreReceipt.ContractAddress).To(Equal(""))
		Expect(coreReceipt.CumulativeGasUsed).To(Equal(expected.CumulativeGasUsed))
		Expect(coreReceipt.GasUsed).To(Equal(expected.GasUsed))
		Expect(coreReceipt.Logs).To(Equal(expected.Logs))
		Expect(coreReceipt.StateRoot).To(Equal(expected.StateRoot))
		Expect(coreReceipt.Status).To(Equal(expected.Status))
		Expect(coreReceipt.TxHash).To(Equal(expected.TxHash))
	})
})
