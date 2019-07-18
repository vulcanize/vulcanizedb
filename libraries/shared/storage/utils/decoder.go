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

package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func Decode(row StorageDiffRow, metadata StorageValueMetadata) (interface{}, error) {
	switch metadata.Type {
	case Uint256:
		return decodeUint256(row.StorageValue.Bytes()), nil
	case Uint48:
		return decodeUint48(row.StorageValue.Bytes()), nil
	case Uint128:
		return decodeUint128(row.StorageValue.Bytes()), nil
	case Address:
		return decodeAddress(row.StorageValue.Bytes()), nil
	case Bytes32:
		return row.StorageValue.Hex(), nil
	case PackedSlot:
		return decodePackedSlot(row.StorageValue.Bytes(), metadata.PackedTypes), nil
	default:
		panic(fmt.Sprintf("can't decode unknown type: %d", metadata.Type))
	}
}

func decodeUint256(raw []byte) string {
	n := big.NewInt(0).SetBytes(raw)
	return n.String()
}

func decodeUint128(raw []byte) string {
	n := big.NewInt(0).SetBytes(raw)
	return n.String()
}

func decodeUint48(raw []byte) string {
	n := big.NewInt(0).SetBytes(raw)
	return n.String()
}

func decodeAddress(raw []byte) string {
	return common.BytesToAddress(raw).Hex()
}

func decodePackedSlot(raw []byte, packedTypes map[int]ValueType) []string{
	storageSlot := raw
	var results []string
	//the reason we're using a map and not a slice is that golang doesn't guarantee the order of a slice
	for _, valueType := range packedTypes {
		lengthOfStorageSlot := len(storageSlot)
		lengthOfItem := getNumberOfBytes(valueType)
		itemStartingIndex := lengthOfStorageSlot - lengthOfItem
		value := storageSlot[itemStartingIndex:]
		decodedValue := decodeIndividualItems(value, valueType)
		results = append(results, decodedValue)

		//pop last item off slot before moving on
		storageSlot = storageSlot[0:itemStartingIndex]
	}

	return results
}

func decodeIndividualItems(itemBytes []byte, valueType ValueType) string {
	switch valueType {
	case Uint48:
		return decodeUint48(itemBytes)
	case Uint128:
		return decodeUint128(itemBytes)
	default:
		panic(fmt.Sprintf("can't decode unknown type: %d", valueType))
	}
}

func getNumberOfBytes(valueType ValueType) int{
	// 8 bits per byte
	switch valueType {
	case Uint48:
		return 48 / 8
	case Uint128:
		return 128 / 8
	default:
		panic(fmt.Sprintf("ValueType %d not recognized", valueType))
	}
}
