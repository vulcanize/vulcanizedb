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

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// PayloadConverter satisfies the PayloadConverter interface for bitcoin
type PayloadConverter struct{}

// NewPayloadConverter creates a pointer to a new PayloadConverter which satisfies the PayloadConverter interface
func NewPayloadConverter() *PayloadConverter {
	return &PayloadConverter{}
}

// Convert method is used to convert a bitcoin BlockPayload to an IPLDPayload
// Satisfies the shared.PayloadConverter interface
func (pc *PayloadConverter) Convert(payload interface{}) (interface{}, error) {
	btcBlockPayload, ok := payload.(BlockPayload)
	if !ok {
		return nil, fmt.Errorf("btc converter: expected payload type %T got %T", BlockPayload{}, payload)
	}
	msgBlock := wire.NewMsgBlock(btcBlockPayload.Header)
	for _, tx := range btcBlockPayload.Txs {
		msgBlock.AddTransaction(tx.MsgTx())
	}
	w := bytes.NewBuffer(make([]byte, 0, msgBlock.SerializeSize()))
	if err := msgBlock.Serialize(w); err != nil {
		return nil, err
	}
	utilBlock := btcutil.NewBlockFromBlockAndBytes(msgBlock, w.Bytes())
	utilBlock.SetHeight(btcBlockPayload.Height)
	txMeta := make([]TxModel, len(btcBlockPayload.Txs))
	for _, tx := range utilBlock.Transactions() {
		index := tx.Index()
		txModel := TxModel{
			TxHash:     tx.Hash().String(),
			Index:      int64(tx.Index()),
			HasWitness: tx.HasWitness(),
		}
		if tx.HasWitness() {
			txModel.WitnessHash = tx.WitnessHash().String()
		}
		txMeta[index] = txModel
	}
	return IPLDPayload{
		Block:      utilBlock,
		TxMetaData: txMeta,
	}, nil
}
