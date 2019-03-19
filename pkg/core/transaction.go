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

package core

type Transaction struct {
	Data     string `db:"input_data"`
	From     string `db:"tx_from"`
	GasLimit uint64
	GasPrice int64
	Hash     string
	Nonce    uint64
	Raw      []byte
	Receipt
	To      string `db:"tx_to"`
	TxIndex int64  `db:"tx_index"`
	Value   string
}
