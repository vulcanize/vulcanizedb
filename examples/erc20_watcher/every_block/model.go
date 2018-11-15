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

package every_block

// Struct to hold token supply data
type TokenSupply struct {
	Value        string
	TokenAddress string
	BlockNumber  int64
}

// Struct to hold token holder address balance data
type TokenBalance struct {
	Value              string
	TokenAddress       string
	BlockNumber        int64
	TokenHolderAddress string
}

// Struct to hold token allowance data
type TokenAllowance struct {
	Value               string
	TokenAddress        string
	BlockNumber         int64
	TokenHolderAddress  string
	TokenSpenderAddress string
}
