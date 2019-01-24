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

package bite

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type BiteEntity struct {
	Ilk              [32]byte
	Urn              [32]byte
	Ink              *big.Int
	Art              *big.Int
	Tab              *big.Int
	Flip             *big.Int
	IArt             *big.Int
	LogIndex         uint
	TransactionIndex uint
	Raw              types.Log
}
