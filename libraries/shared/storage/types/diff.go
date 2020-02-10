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

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
)

const ExpectedRowLength = 5

type RawDiff struct {
	HashedAddress common.Hash `db:"hashed_address"`
	BlockHash     common.Hash `db:"block_hash"`
	BlockHeight   int         `db:"block_height"`
	StorageKey    common.Hash `db:"storage_key"`
	StorageValue  common.Hash `db:"storage_value"`
}

type PersistedDiff struct {
	RawDiff
	Checked  bool
	ID       int64
	HeaderID int64 `db:"header_id"`
}

func FromParityCsvRow(csvRow []string) (RawDiff, error) {
	if len(csvRow) != ExpectedRowLength {
		return RawDiff{}, ErrRowMalformed{Length: len(csvRow)}
	}
	height, err := strconv.Atoi(csvRow[2])
	if err != nil {
		return RawDiff{}, err
	}
	return RawDiff{
		HashedAddress: HexToKeccak256Hash(csvRow[0]),
		BlockHash:     common.HexToHash(csvRow[1]),
		BlockHeight:   height,
		StorageKey:    common.HexToHash(csvRow[3]),
		StorageValue:  common.HexToHash(csvRow[4]),
	}, nil
}

func FromGethStateDiff(account statediff.AccountDiff, stateDiff *statediff.StateDiff, storage statediff.StorageDiff) (RawDiff, error) {
	var decodedRLPStorageValue []byte
	err := rlp.DecodeBytes(storage.Value, &decodedRLPStorageValue)
	if err != nil {
		return RawDiff{}, err
	}

	return RawDiff{
		HashedAddress: common.BytesToHash(account.Key),
		BlockHash:     stateDiff.BlockHash,
		BlockHeight:   int(stateDiff.BlockNumber.Int64()),
		StorageKey:    common.BytesToHash(storage.Key),
		StorageValue:  common.BytesToHash(decodedRLPStorageValue),
	}, nil
}

func ToPersistedDiff(raw RawDiff, id int64) PersistedDiff {
	return PersistedDiff{
		RawDiff: raw,
		ID:      id,
	}
}

func HexToKeccak256Hash(hex string) common.Hash {
	return crypto.Keccak256Hash(common.FromHex(hex))
}
