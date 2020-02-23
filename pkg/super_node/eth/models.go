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

// HeaderModel is the db model for eth.header_cids
type HeaderModel struct {
	ID              int64  `db:"id"`
	BlockNumber     string `db:"block_number"`
	BlockHash       string `db:"block_hash"`
	ParentHash      string `db:"parent_hash"`
	CID             string `db:"cid"`
	TotalDifficulty string `db:"td"`
	NodeID          int64  `db:"node_id"`
}

// UncleModel is the db model for eth.uncle_cids
type UncleModel struct {
	ID         int64  `db:"id"`
	HeaderID   int64  `db:"header_id"`
	BlockHash  string `db:"block_hash"`
	ParentHash string `db:"parent_hash"`
	CID        string `db:"cid"`
}

// TxModel is the db model for eth.transaction_cids
type TxModel struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	Index    int64  `db:"index"`
	TxHash   string `db:"tx_hash"`
	CID      string `db:"cid"`
	Dst      string `db:"dst"`
	Src      string `db:"src"`
}

// ReceiptModel is the db model for eth.receipt_cids
type ReceiptModel struct {
	ID       int64          `db:"id"`
	TxID     int64          `db:"tx_id"`
	CID      string         `db:"cid"`
	Contract string         `db:"contract"`
	Topic0s  pq.StringArray `db:"topic0s"`
	Topic1s  pq.StringArray `db:"topic1s"`
	Topic2s  pq.StringArray `db:"topic2s"`
	Topic3s  pq.StringArray `db:"topic3s"`
}

// StateNodeModel is the db model for eth.state_cids
type StateNodeModel struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	StateKey string `db:"state_key"`
	Leaf     bool   `db:"leaf"`
	CID      string `db:"cid"`
}

// StorageNodeModel is the db model for eth.storage_cids
type StorageNodeModel struct {
	ID         int64  `db:"id"`
	StateID    int64  `db:"state_id"`
	StorageKey string `db:"storage_key"`
	Leaf       bool   `db:"leaf"`
	CID        string `db:"cid"`
}

// StorageNodeWithStateKeyModel is a db model for eth.storage_cids + eth.state_cids.state_key
type StorageNodeWithStateKeyModel struct {
	ID         int64  `db:"id"`
	StateID    int64  `db:"state_id"`
	StateKey   string `db:"state_key"`
	StorageKey string `db:"storage_key"`
	Leaf       bool   `db:"leaf"`
	CID        string `db:"cid"`
}
