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

import "github.com/lib/pq"

// HeaderModel is the db model for btc.header_cids table
type HeaderModel struct {
	ID          int64  `db:"id"`
	BlockNumber string `db:"block_number"`
	BlockHash   string `db:"block_hash"`
	ParentHash  string `db:"parent_hash"`
	CID         string `db:"cid"`
	Version     int32  `db:"version"`
	Timestamp   int64  `db:"timestamp"`
	Bits        uint32 `db:"bits"`
}

// TxModel is the db model for btc.transaction_cids table
type TxModel struct {
	ID          int64  `db:"id"`
	HeaderID    int64  `db:"header_id"`
	Index       int64  `db:"index"`
	TxHash      string `db:"tx_hash"`
	CID         string `db:"cid"`
	SegWit      bool   `db:"segwit"`
	WitnessHash string `db:"witness_hash"`
}

// TxModelWithInsAndOuts is the db model for btc.transaction_cids table that includes the children tx_input and tx_output tables
type TxModelWithInsAndOuts struct {
	ID          int64  `db:"id"`
	HeaderID    int64  `db:"header_id"`
	Index       int64  `db:"index"`
	TxHash      string `db:"tx_hash"`
	CID         string `db:"cid"`
	SegWit      bool   `db:"segwit"`
	WitnessHash string `db:"witness_hash"`
	TxInputs    []TxInput
	TxOutputs   []TxOutput
}

// TxInput is the db model for btc.tx_inputs table
type TxInput struct {
	ID                    int64    `db:"id"`
	TxID                  int64    `db:"tx_id"`
	Index                 int64    `db:"index"`
	TxWitness             [][]byte `db:"witness"`
	SignatureScript       []byte   `db:"sig_script"`
	PreviousOutPointTxID  int64    `db:"outpoint_tx_id"`
	PreviousOutPointIndex uint32   `db:"outpoint_index"`
	PreviousOutPointHash  string
}

// TxOutput is the db model for btc.tx_outputs table
type TxOutput struct {
	ID           int64          `db:"id"`
	TxID         int64          `db:"tx_id"`
	Index        int64          `db:"index"`
	Value        int64          `db:"value"`
	PkScript     []byte         `db:"pk_script"`
	ScriptClass  uint8          `db:"script_class"`
	RequiredSigs int64          `db:"required_sigs"`
	Addresses    pq.StringArray `db:"addresses"`
}
