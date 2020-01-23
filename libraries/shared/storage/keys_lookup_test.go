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
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage keys lookup utils", func() {
	Describe("AddHashedKeys", func() {
		It("returns a copy of the map with an additional slot for the hashed version of every key", func() {
			fakeMap := map[common.Hash]types.ValueMetadata{}
			fakeStorageKey := common.HexToHash("72c72de6b203d67cb6cd54fc93300109fcc6fd6eac88e390271a3d548794d800")
			var fakeMappingKey types.Key = "fakeKey"
			fakeMetadata := types.ValueMetadata{
				Name: "fakeName",
				Keys: map[types.Key]string{fakeMappingKey: "fakeValue"},
				Type: types.Uint48,
			}
			fakeMap[fakeStorageKey] = fakeMetadata

			result := storage.AddHashedKeys(fakeMap)

			Expect(len(result)).To(Equal(2))
			expectedHashedStorageKey := common.HexToHash("2165edb4e1c37b99b60fa510d84f939dd35d5cd1d1c8f299d6456ea09df65a76")
			Expect(fakeMap[fakeStorageKey]).To(Equal(fakeMetadata))
			Expect(fakeMap[expectedHashedStorageKey]).To(Equal(fakeMetadata))
		})
	})
})
