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

package utils_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

var _ = Describe("Storage decoder", func() {
	It("decodes uint256", func() {
		fakeInt := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000539")
		row := utils.StorageDiffRow{StorageValue: fakeInt}
		metadata := utils.StorageValueMetadata{Type: utils.Uint256}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(big.NewInt(0).SetBytes(fakeInt.Bytes()).String()))
	})

	It("decodes uint48", func() {
		fakeInt := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000123")
		row := utils.StorageDiffRow{StorageValue: fakeInt}
		metadata := utils.StorageValueMetadata{Type: utils.Uint48}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(big.NewInt(0).SetBytes(fakeInt.Bytes()).String()))
	})

	It("decodes address", func() {
		fakeAddress := common.HexToAddress("0x12345")
		row := utils.StorageDiffRow{StorageValue: fakeAddress.Hash()}
		metadata := utils.StorageValueMetadata{Type: utils.Address}

		result, err := utils.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(fakeAddress.Hex()))
	})
})
