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

package helpers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// BigFromString creates a big.Int from a string
func BigFromString(n string) *big.Int {
	b := new(big.Int)
	b.SetString(n, 10)
	return b
}

// GenerateSignature returns the keccak256 hash hex of a string
func GenerateSignature(s string) string {
	eventSignature := []byte(s)
	hash := crypto.Keccak256Hash(eventSignature)
	return hash.Hex()
}
