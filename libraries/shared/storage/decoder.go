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

package storage

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

const (
	bitsPerByte = 8
)

func Decode(diff types.PersistedDiff, metadata types.ValueMetadata) interface{} {
	switch metadata.Type {
	case types.Uint256:
		return decodeInteger(diff.StorageValue.Bytes())
	case types.Uint8:
		return decodeInteger(diff.StorageValue.Bytes())
	case types.Uint32:
		return decodeInteger(diff.StorageValue.Bytes())
	case types.Uint48:
		return decodeInteger(diff.StorageValue.Bytes())
	case types.Uint128:
		return decodeInteger(diff.StorageValue.Bytes())
	case types.Address:
		return decodeAddress(diff.StorageValue.Bytes())
	case types.Bytes32:
		return diff.StorageValue.Hex()
	case types.PackedSlot:
		return decodePackedSlot(diff.StorageValue.Bytes(), metadata.PackedTypes)
	default:
		panic(fmt.Sprintf("can't decode unknown type: %d", metadata.Type))
	}
}

func decodeInteger(raw []byte) string {
	n := big.NewInt(0).SetBytes(raw)
	return n.String()
}

func decodeAddress(raw []byte) string {
	return common.BytesToAddress(raw).Hex()
}

func decodePackedSlot(raw []byte, packedTypes map[int]types.ValueType) map[int]string {
	storageSlotData := raw
	decodedStorageSlotItems := map[int]string{}
	numberOfTypes := len(packedTypes)

	for position := 0; position < numberOfTypes; position++ {
		//get length of remaining storage date
		lengthOfStorageData := len(storageSlotData)

		//get item details (type, length, starting index, value bytes)
		itemType := packedTypes[position]
		lengthOfItem := getNumberOfBytes(itemType)
		itemStartingIndex := lengthOfStorageData - lengthOfItem
		itemValueBytes := storageSlotData[itemStartingIndex:]

		//decode item's bytes and set in results map
		decodedValue := decodeIndividualItem(itemValueBytes, itemType)
		decodedStorageSlotItems[position] = decodedValue

		//pop last item off raw slot data before moving on
		storageSlotData = storageSlotData[0:itemStartingIndex]
	}

	return decodedStorageSlotItems
}

func decodeIndividualItem(itemBytes []byte, valueType types.ValueType) string {
	switch valueType {
	case types.Uint32, types.Uint48, types.Uint128:
		return decodeInteger(itemBytes)
	case types.Address:
		return decodeAddress(itemBytes)
	default:
		panic(fmt.Sprintf("can't decode unknown type: %d", valueType))
	}
}

func getNumberOfBytes(valueType types.ValueType) int {
	switch valueType {
	case types.Uint32:
		return 32 / bitsPerByte
	case types.Uint48:
		return 48 / bitsPerByte
	case types.Uint128:
		return 128 / bitsPerByte
	case types.Address:
		return 20
	default:
		panic(fmt.Sprintf("ValueType %d not recognized", valueType))
	}
}
