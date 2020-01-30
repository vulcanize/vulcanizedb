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

package btc

type HeaderModel struct {
	ID          int64  `db:"id"`
	BlockNumber string `db:"block_number"`
	BlockHash   string `db:"block_hash"`
	ParentHash  string `db:"parent_hash"`
	CID         string `db:"cid"`
	// TotalDifficulty string `db:"td"` Not sure we can support this unless we modify btcd
}

type TxModel struct {
	ID          int64  `db:"id"`
	HeaderID    int64  `db:"header_id"`
	Index       int64  `db:"index"`
	TxHash      string `db:"tx_hash"`
	CID         string `db:"cid"`
	HasWitness  bool   `db:"has_witness"`
	WitnessHash string `db:"witness_hash"`
	//Dst      string `db:"dst"`
	//Src      string `db:"src"`
}
