package storage_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"math/big"
)

var _ = Describe("Mappings", func() {
	Describe("GetMapping", func() {
		It("returns the storage key for a mapping when passed the mapping's index on the contract and the desired value's key", func() {
			// ex. solidity:
			//    	mapping (bytes32 => uint) public amounts
			// pass in the index of the mapping on the contract + the bytes32 key for the uint val being looked up
			indexOfMappingOnContract := storage.IndexZero
			keyForDesiredValueInMapping := "fake_bytes32"

			storageKey := storage.GetMapping(indexOfMappingOnContract, keyForDesiredValueInMapping)

			expectedStorageKeyBytes := crypto.Keccak256(common.FromHex(keyForDesiredValueInMapping + indexOfMappingOnContract))
			expectedStorageKey := common.BytesToHash(expectedStorageKeyBytes)
			Expect(storageKey).To(Equal(expectedStorageKey))
		})
	})

	Describe("GetNestedMapping", func() {
		It("returns the storage key for a nested mapping when passed the mapping's index on the contract and the desired value's keys", func() {
			// ex. solidity:
			//    	mapping (address => mapping (uint => bytes32)) public addressNames
			// pass in the index of the mapping on the contract + the address and uint keys for the bytes32 val being looked up
			indexOfMappingOnContract := storage.IndexOne
			keyForOuterMapping := "fake_address"
			keyForInnerMapping := "123"

			storageKey := storage.GetNestedMapping(indexOfMappingOnContract, keyForOuterMapping, keyForInnerMapping)

			hashedOuterMappingStorageKey := crypto.Keccak256(common.FromHex(keyForOuterMapping + indexOfMappingOnContract))
			fullStorageKeyBytes := crypto.Keccak256(common.FromHex(keyForInnerMapping), hashedOuterMappingStorageKey)
			expectedStorageKey := common.BytesToHash(fullStorageKeyBytes)
			Expect(storageKey).To(Equal(expectedStorageKey))
		})
	})

	Describe("GetIncrementedKey", func() {
		It("returns the storage key for later values sharing an index on the contract with other earlier values", func() {
			// ex. solidity:
			//    	struct Data {
			//        uint256 quantity;
			//        uint256 quality;
			//    	}
			//    	mapping (bytes32 => Data) public itemData;
			// pass in the storage key for the zero-indexed value ("quantity") + the number of increments required.
			// (For "quality", we must increment the storage key for the corresponding "quantity" by 1).
			indexOfMappingOnContract := storage.IndexTwo
			keyForDesiredValueInMapping := "fake_bytes32"
			storageKeyForFirstPropertyOnStruct := storage.GetMapping(indexOfMappingOnContract, keyForDesiredValueInMapping)

			storageKey := storage.GetIncrementedKey(storageKeyForFirstPropertyOnStruct, 1)

			incrementedStorageKey := storageKeyForFirstPropertyOnStruct.Big().Add(storageKeyForFirstPropertyOnStruct.Big(), big.NewInt(1))
			expectedStorageKey := common.BytesToHash(incrementedStorageKey.Bytes())
			Expect(storageKey).To(Equal(expectedStorageKey))
		})
	})
})
