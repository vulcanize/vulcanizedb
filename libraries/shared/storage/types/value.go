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

package types

import "fmt"

type ValueType int

const (
	Uint256 ValueType = iota
	Uint32
	Uint48
	Uint128
	Bytes32
	Address
	PackedSlot
)

type Key string

type ValueMetadata struct {
	Name        string
	Keys        map[Key]string
	Type        ValueType
	PackedNames map[int]string    //zero indexed position in map => name of packed item
	PackedTypes map[int]ValueType //zero indexed position in map => type of packed item
}

func GetValueMetadata(name string, keys map[Key]string, valueType ValueType) ValueMetadata {
	return getMetadata(name, keys, valueType, nil, nil)
}

func GetValueMetadataForPackedSlot(name string, keys map[Key]string, valueType ValueType, packedNames map[int]string, packedTypes map[int]ValueType) ValueMetadata {
	return getMetadata(name, keys, valueType, packedNames, packedTypes)
}

func getMetadata(name string, keys map[Key]string, valueType ValueType, packedNames map[int]string, packedTypes map[int]ValueType) ValueMetadata {
	assertPackedSlotArgs(valueType, packedNames, packedTypes)

	return ValueMetadata{
		Name:        name,
		Keys:        keys,
		Type:        valueType,
		PackedNames: packedNames,
		PackedTypes: packedTypes,
	}
}

func assertPackedSlotArgs(valueType ValueType, packedNames map[int]string, packedTypes map[int]ValueType) {
	if valueType == PackedSlot && (packedTypes == nil || packedNames == nil) {
		panic(fmt.Sprintf("ValueType is PackedSlot. Expected PackedNames and PackedTypes to not be nil, but got PackedNames = %v and PackedTypes = %v", packedNames, packedTypes))
	} else if (packedNames != nil && packedTypes != nil) && valueType != PackedSlot {
		panic(fmt.Sprintf("PackedNames and PackedTypes passed in. Expected ValueType to equal PackedSlot (%v), but got %v.", PackedSlot, valueType))
	}

}
