// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	case Address:
		return decodeAddress(row.StorageValue.Bytes()), nil
	case Bytes32:
		return row.StorageValue.Hex(), nil
	default:
		panic(fmt.Sprintf("can't decode unknown type: %d", metadata.Type))
	}
}

func decodeUint256(raw []byte) string {
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
