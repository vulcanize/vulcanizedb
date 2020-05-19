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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
)

const ExpectedRowLength = 5

type StorageDiffInput struct {
	HashedAddress common.Hash `db:"hashed_address"`
	BlockHash     common.Hash `db:"block_hash"`
	BlockHeight   int         `db:"block_height"`
	StorageKey    common.Hash `db:"storage_key"`
	StorageValue  common.Hash `db:"storage_value"`
}

type PersistedStorageDiff struct {
	StorageDiffInput
	ID int64
}

func FromParityCsvRow(csvRow []string) (StorageDiffInput, error) {
	if len(csvRow) != ExpectedRowLength {
		return StorageDiffInput{}, ErrRowMalformed{Length: len(csvRow)}
	}
	height, err := strconv.Atoi(csvRow[2])
	if err != nil {
		return StorageDiffInput{}, err
	}
	return StorageDiffInput{
		HashedAddress: HexToKeccak256Hash(csvRow[0]),
		BlockHash:     common.HexToHash(csvRow[1]),
		BlockHeight:   height,
		StorageKey:    common.HexToHash(csvRow[3]),
		StorageValue:  common.HexToHash(csvRow[4]),
	}, nil
}

func FromGethStateDiff(account statediff.StateNode, stateDiff *statediff.StateObject, storage statediff.StorageNode) (StorageDiffInput, error) {
	var decodedValue []byte
	err := rlp.DecodeBytes(storage.NodeValue, &decodedValue)
	if err != nil {
		return StorageDiffInput{}, err
	}

	return StorageDiffInput{
		HashedAddress: common.BytesToHash(account.LeafKey),
		BlockHash:     stateDiff.BlockHash,
		BlockHeight:   int(stateDiff.BlockNumber.Int64()),
		StorageKey:    common.BytesToHash(storage.LeafKey),
		StorageValue:  common.BytesToHash(decodedValue),
	}, nil
}

func ToPersistedDiff(raw StorageDiffInput, id int64) PersistedStorageDiff {
	return PersistedStorageDiff{
		StorageDiffInput: raw,
		ID:               id,
	}
}

func HexToKeccak256Hash(hex string) common.Hash {
	return crypto.Keccak256Hash(common.FromHex(hex))
}

func GetAccountsFromDiff(stateDiff statediff.StateObject) []statediff.StateNode {
	return stateDiff.Nodes
}
