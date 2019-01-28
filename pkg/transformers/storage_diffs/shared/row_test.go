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
)

var _ = Describe("Storage row parsing", func() {
	It("converts an array of strings to a row struct", func() {
		contract := "0x123"
		blockHash := "0x456"
		blockHeight := "789"
		storageKey := "0x987"
		storageValue := "0x654"
		data := []string{contract, blockHash, blockHeight, storageKey, storageValue}

		result, err := shared.FromStrings(data)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Contract).To(Equal(common.HexToAddress(contract)))
		Expect(result.BlockHash).To(Equal(common.HexToHash(blockHash)))
		Expect(result.BlockHeight).To(Equal(789))
		Expect(result.StorageKey).To(Equal(common.HexToHash(storageKey)))
		Expect(result.StorageValue).To(Equal(common.HexToHash(storageValue)))
	})

	It("returns an error if row is missing data", func() {
		_, err := shared.FromStrings([]string{"0x123"})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(shared.ErrRowMalformed{Length: 1}))
	})

	It("returns error if block height malformed", func() {
		_, err := shared.FromStrings([]string{"", "", "", "", ""})

		Expect(err).To(HaveOccurred())
	})
})
