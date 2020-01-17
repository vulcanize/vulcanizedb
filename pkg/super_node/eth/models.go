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

package eth

import "github.com/lib/pq"

type HeaderModel struct {
	ID              int64  `db:"id"`
	BlockNumber     string `db:"block_number"`
	BlockHash       string `db:"block_hash"`
	ParentHash      string `db:"parent_hash"`
	CID             string `db:"cid"`
	Uncle           bool   `db:"uncle"`
	TotalDifficulty string `db:"td"`
}

type TxModel struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	TxHash   string `db:"tx_hash"`
	CID      string `db:"cid"`
	Dst      string `db:"dst"`
	Src      string `db:"src"`
}

type ReceiptModel struct {
	ID       int64          `db:"id"`
	TxID     int64          `db:"tx_id"`
	CID      string         `db:"cid"`
	Contract string         `db:"contract"`
	Topic0s  pq.StringArray `db:"topic0s"`
}

type StateNodeModel struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	StateKey string `db:"state_key"`
	Leaf     bool   `db:"leaf"`
	CID      string `db:"cid"`
}

type StorageNodeModel struct {
	ID         int64  `db:"id"`
	StateID    int64  `db:"state_id"`
	StorageKey string `db:"storage_key"`
	Leaf       bool   `db:"leaf"`
	CID        string `db:"cid"`
}

type StorageNodeWithStateKeyModel struct {
	ID         int64  `db:"id"`
	StateID    int64  `db:"state_id"`
	StateKey   string `db:"state_key"`
	StorageKey string `db:"storage_key"`
	Leaf       bool   `db:"leaf"`
	CID        string `db:"cid"`
}
