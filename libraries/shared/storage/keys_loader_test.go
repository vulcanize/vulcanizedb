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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage keys loader utils", func() {
	Describe("GetStorageKeyForMapping", func() {
		It("returns the storage key for a mapping when passed the mapping's index on the contract and the desired value's key", func() {
			// ex. solidity:
			//    	mapping (bytes32 => uint) public amounts
			// to access amounts, pass in the index of the mapping on the contract + the bytes32 key for the uint val being looked up
			indexOfMappingOnContract := storage.IndexZero
			keyForDesiredValueInMapping := "1234567890abcdef"

			storageKey := storage.GetStorageKeyForMapping(indexOfMappingOnContract, keyForDesiredValueInMapping)

			expectedStorageKey := common.HexToHash("0xee0c1b59a3856bafbfb8730e7694c4badc271eb5f01ce4a8d7a53d8a6499676f")
			Expect(storageKey).To(Equal(expectedStorageKey))
		})

		It("returns same result if value includes hex prefix", func() {
			indexOfMappingOnContract := storage.IndexZero
			keyForDesiredValueInMapping := "0x1234567890abcdef"

			storageKey := storage.GetStorageKeyForMapping(indexOfMappingOnContract, keyForDesiredValueInMapping)

			expectedStorageKey := common.HexToHash("0xee0c1b59a3856bafbfb8730e7694c4badc271eb5f01ce4a8d7a53d8a6499676f")
			Expect(storageKey).To(Equal(expectedStorageKey))
		})
	})

	Describe("GetStorageKeyForNestedMapping", func() {
		It("returns the storage key for a nested mapping when passed the mapping's index on the contract and the desired value's keys", func() {
			// ex. solidity:
			//    	mapping (bytes32 => uint) public amounts
			//    	mapping (address => mapping (uint => bytes32)) public addressNames
			// to access addressNames, pass in the index of the mapping on the contract + the address and uint keys for the bytes32 val being looked up
			indexOfMappingOnContract := storage.IndexOne
			keyForOuterMapping := "1234567890abcdef"
			keyForInnerMapping := "123"

			storageKey := storage.GetStorageKeyForNestedMapping(indexOfMappingOnContract, keyForOuterMapping, keyForInnerMapping)

			expectedStorageKey := common.HexToHash("0x82113529f6cd61061d1a6f0de53f2bdd067a1addd3d2b46be50a99abfcdb1661")
			Expect(storageKey).To(Equal(expectedStorageKey))
		})
	})

	Describe("GetIncrementedStorageKey", func() {
		It("returns the storage key for later values sharing an index on the contract with other earlier values", func() {
			// ex. solidity:
			//    	mapping (bytes32 => uint) public amounts
			//    	mapping (address => mapping (uint => bytes32)) public addressNames
			//    	struct Data {
			//        uint256 quantity;
			//        uint256 quality;
			//    	}
			//    	mapping (bytes32 => Data) public itemData;
			// to access quality from itemData, pass in the storage key for the zero-indexed value (quantity) + the number of increments required.
			// (For "quality", we must increment the storage key for the corresponding "quantity" by 1).
			indexOfMappingOnContract := storage.IndexTwo
			keyForDesiredValueInMapping := "1234567890abcdef"
			storageKeyForFirstPropertyOnStruct := storage.GetStorageKeyForMapping(indexOfMappingOnContract, keyForDesiredValueInMapping)

			storageKey := storage.GetIncrementedStorageKey(storageKeyForFirstPropertyOnStruct, 1)

			expectedStorageKey := common.HexToHash("0x69b38749f0a8ed5d505c8474f7fb62c7828aad8a7627f1c67e07af1d2368cad4")
			Expect(storageKey).To(Equal(expectedStorageKey))
		})
	})
})
