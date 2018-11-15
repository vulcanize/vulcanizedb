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

package dai

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransferEntity struct {
	TokenName    string
	TokenAddress common.Address
	Src          common.Address
	Dst          common.Address
	Wad          *big.Int
	Block        int64
	TxHash       string
}

type ApprovalEntity struct {
	TokenName    string
	TokenAddress common.Address
	Src          common.Address
	Guy          common.Address
	Wad          *big.Int
	Block        int64
	TxHash       string
}
