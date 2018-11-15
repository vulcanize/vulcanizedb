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

package event_triggered

type BurnModel struct {
	TokenName    string `db:"token_name"`
	TokenAddress string `db:"token_address"`
	Burner       string `db:"burner"`
	Tokens       string `db:"tokens"`
	Block        int64  `db:"block"`
	TxHash       string `db:"tx"`
}

type MintModel struct {
	TokenName    string `db:"token_name"`
	TokenAddress string `db:"token_address"`
	Mintee       string `db:"mintee"`
	Minter       string `db:"minter"`
	Tokens       string `db:"tokens"`
	Block        int64  `db:"block"`
	TxHash       string `db:"tx"`
}
