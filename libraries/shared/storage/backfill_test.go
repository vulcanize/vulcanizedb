// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package storage_test

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
)

var _ = Describe("BackFiller", func() {
	Describe("BackFill", func() {
		var (
			fetcher    *mocks.StateDiffFetcher
			backFiller storage.BackFiller
		)
		BeforeEach(func() {
			fetcher = new(mocks.StateDiffFetcher)
			fetcher.SetPayloadsToReturn(map[uint64]*statediff.Payload{
				test_data.BlockNumber.Uint64():  &test_data.MockStatediffPayload,
				test_data.BlockNumber2.Uint64(): &test_data.MockStatediffPayload2,
			})
			backFiller = storage.NewStorageBackFiller(fetcher)
		})
		It("Batch calls statediff_stateDiffAt", func() {
			backFillStorage, backFillErr := backFiller.BackFill(test_data.BlockNumber.Uint64(), test_data.BlockNumber2.Uint64())
			Expect(backFillErr).ToNot(HaveOccurred())
			Expect(len(backFillStorage)).To(Equal(4))
			// Can only rlp encode the slice of diffs as part of a struct
			// Rlp encoding allows us to compare content of the slices when the order in the slice may vary
			expectedDiffStruct := struct {
				diffs []utils.StorageDiff
			}{
				[]utils.StorageDiff{
					test_data.CreatedExpectedStorageDiff,
					test_data.UpdatedExpectedStorageDiff,
					test_data.UpdatedExpectedStorageDiff2,
					test_data.DeletedExpectedStorageDiff,
				},
			}
			expectedDiffBytes, rlpErr1 := rlp.EncodeToBytes(expectedDiffStruct)
			Expect(rlpErr1).ToNot(HaveOccurred())
			receivedDiffStruct := struct {
				diffs []utils.StorageDiff
			}{
				backFillStorage,
			}
			receivedDiffBytes, rlpErr2 := rlp.EncodeToBytes(receivedDiffStruct)
			Expect(rlpErr2).ToNot(HaveOccurred())
			Expect(bytes.Equal(expectedDiffBytes, receivedDiffBytes)).To(BeTrue())
		})
	})
})
