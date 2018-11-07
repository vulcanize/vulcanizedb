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
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	common2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
	"math/big"
)

var _ = Describe("Block header converter", func() {
	It("converts geth header to core header", func() {
		gethHeader := &types.Header{
			Difficulty:  big.NewInt(1),
			Number:      big.NewInt(2),
			ParentHash:  common.HexToHash("0xParent"),
			ReceiptHash: common.HexToHash("0xReceipt"),
			Root:        common.HexToHash("0xRoot"),
			Time:        big.NewInt(3),
			TxHash:      common.HexToHash("0xTransaction"),
			UncleHash:   common.HexToHash("0xUncle"),
		}
		converter := common2.HeaderConverter{}

		coreHeader, err := converter.Convert(gethHeader)

		Expect(err).NotTo(HaveOccurred())
		Expect(coreHeader.BlockNumber).To(Equal(gethHeader.Number.Int64()))
		Expect(coreHeader.Hash).To(Equal(gethHeader.Hash().Hex()))
	})

	It("includes raw bytes for header", func() {
		headerRLP := []byte{249, 2, 23, 160, 180, 251, 173, 248, 234, 69, 43, 19, 151, 24, 226, 112, 13, 193, 19, 92, 252, 129, 20, 80, 49, 200, 75, 122, 178, 124, 215, 16, 57, 79, 123, 56, 160, 29, 204, 77, 232, 222, 199, 93, 122, 171, 133, 181, 103, 182, 204, 212, 26, 211, 18, 69, 27, 148, 138, 116, 19, 240, 161, 66, 253, 64, 212, 147, 71, 148, 42, 101, 172, 164, 213, 252, 91, 92, 133, 144, 144, 166, 195, 77, 22, 65, 53, 57, 130, 38, 160, 14, 6, 111, 60, 34, 151, 165, 203, 48, 5, 147, 5, 38, 23, 209, 188, 165, 148, 111, 12, 170, 6, 53, 253, 177, 184, 90, 199, 229, 35, 111, 52, 160, 101, 186, 136, 127, 203, 8, 38, 246, 22, 208, 31, 115, 108, 29, 45, 103, 123, 202, 189, 226, 247, 252, 37, 170, 145, 207, 188, 11, 59, 173, 92, 179, 160, 32, 227, 83, 69, 64, 202, 241, 99, 120, 230, 232, 106, 43, 241, 35, 109, 159, 135, 109, 50, 24, 251, 192, 57, 88, 230, 219, 28, 99, 75, 35, 51, 185, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 134, 11, 105, 222, 129, 162, 43, 131, 15, 66, 64, 131, 47, 239, 216, 130, 196, 68, 132, 86, 191, 180, 21, 152, 215, 131, 1, 3, 3, 132, 71, 101, 116, 104, 135, 103, 111, 49, 46, 53, 46, 49, 133, 108, 105, 110, 117, 120, 160, 146, 196, 18, 154, 10, 226, 54, 27, 69, 42, 158, 222, 236, 229, 92, 18, 236, 238, 171, 134, 99, 22, 25, 94, 61, 135, 252, 27, 0, 91, 102, 69, 136, 205, 76, 85, 185, 65, 207, 144, 21}
		var gethHeader types.Header
		err := rlp.Decode(bytes.NewReader(headerRLP), &gethHeader)
		Expect(err).NotTo(HaveOccurred())
		converter := common2.HeaderConverter{}

		coreHeader, err := converter.Convert(&gethHeader)

		Expect(err).NotTo(HaveOccurred())
		Expect(coreHeader.Raw).To(Equal(headerRLP))
	})
})
