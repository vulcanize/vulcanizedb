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

package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func Decode(row StorageDiffRow, metadata StorageValueMetadata) (interface{}, error) {
	switch metadata.Type {
	case Uint256:
		return decodeUint256(row.StorageValue.Bytes()), nil
	case Address:
		return decodeAddress(row.StorageValue.Bytes()), nil
	case Bytes32:
		return row.StorageValue.Hex(), nil
	default:
		return nil, ErrTypeNotFound{}
	}
}

func decodeUint256(raw []byte) string {
	n := big.NewInt(0).SetBytes(raw)
	return n.String()
}

func decodeAddress(raw []byte) string {
	return common.BytesToAddress(raw).Hex()
}
