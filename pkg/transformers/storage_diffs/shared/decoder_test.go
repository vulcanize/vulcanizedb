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

package shared_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"math/big"
)

var _ = Describe("Storage decoder", func() {
	It("decodes uint256", func() {
		fakeInt := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000539")
		row := shared.StorageDiffRow{StorageValue: fakeInt}
		metadata := shared.StorageValueMetadata{Type: shared.Uint256}

		result, err := shared.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(big.NewInt(0).SetBytes(fakeInt.Bytes()).String()))
	})

	It("decodes address", func() {
		fakeAddress := common.HexToAddress("0x12345")
		row := shared.StorageDiffRow{StorageValue: fakeAddress.Hash()}
		metadata := shared.StorageValueMetadata{Type: shared.Address}

		result, err := shared.Decode(row, metadata)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(fakeAddress.Hex()))
	})
})
