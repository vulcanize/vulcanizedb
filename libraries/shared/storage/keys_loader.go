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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	IndexZero   = "0000000000000000000000000000000000000000000000000000000000000000"
	IndexOne    = "0000000000000000000000000000000000000000000000000000000000000001"
	IndexTwo    = "0000000000000000000000000000000000000000000000000000000000000002"
	IndexThree  = "0000000000000000000000000000000000000000000000000000000000000003"
	IndexFour   = "0000000000000000000000000000000000000000000000000000000000000004"
	IndexFive   = "0000000000000000000000000000000000000000000000000000000000000005"
	IndexSix    = "0000000000000000000000000000000000000000000000000000000000000006"
	IndexSeven  = "0000000000000000000000000000000000000000000000000000000000000007"
	IndexEight  = "0000000000000000000000000000000000000000000000000000000000000008"
	IndexNine   = "0000000000000000000000000000000000000000000000000000000000000009"
	IndexTen    = "000000000000000000000000000000000000000000000000000000000000000a"
	IndexEleven = "000000000000000000000000000000000000000000000000000000000000000b"
)

func GetKeyForMapping(indexOnContract, key string) common.Hash {
	keyBytes := common.FromHex(key + indexOnContract)
	return crypto.Keccak256Hash(keyBytes)
}

func GetKeyForNestedMapping(indexOnContract, primaryKey, secondaryKey string) common.Hash {
	primaryMappingIndex := crypto.Keccak256(common.FromHex(primaryKey + indexOnContract))
	return crypto.Keccak256Hash(common.FromHex(secondaryKey), primaryMappingIndex)
}

func GetIncrementedKey(original common.Hash, incrementBy int64) common.Hash {
	originalMappingAsInt := original.Big()
	incremented := big.NewInt(0).Add(originalMappingAsInt, big.NewInt(incrementBy))
	return common.BytesToHash(incremented.Bytes())
}
