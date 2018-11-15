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
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	common2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

var _ = Describe("Block header converter", func() {
	It("converts geth header to core header", func() {
		gethHeader := &types.Header{
			Difficulty:  big.NewInt(1),
			Number:      big.NewInt(2),
			ParentHash:  common.HexToHash("0xParent"),
			ReceiptHash: common.HexToHash("0xReceipt"),
			Root:        common.HexToHash("0xRoot"),
			Time:        big.NewInt(123456789),
			TxHash:      common.HexToHash("0xTransaction"),
			UncleHash:   common.HexToHash("0xUncle"),
		}
		converter := common2.HeaderConverter{}
		hash := fakes.FakeHash.String()

		coreHeader, err := converter.Convert(gethHeader, hash)

		Expect(err).NotTo(HaveOccurred())
		Expect(coreHeader.BlockNumber).To(Equal(gethHeader.Number.Int64()))
		Expect(coreHeader.Hash).To(Equal(hash))
		Expect(coreHeader.Timestamp).To(Equal(gethHeader.Time.String()))
	})

	It("includes raw bytes for header as JSON", func() {
		gethHeader := types.Header{Number: big.NewInt(123)}
		converter := common2.HeaderConverter{}

		coreHeader, err := converter.Convert(&gethHeader, fakes.FakeHash.String())

		Expect(err).NotTo(HaveOccurred())
		expectedJSON, err := json.Marshal(gethHeader)
		Expect(err).NotTo(HaveOccurred())
		Expect(coreHeader.Raw).To(Equal(expectedJSON))
	})
})
