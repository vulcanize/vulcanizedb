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
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// PayloadConverter satisfies the PayloadConverter interface for bitcoin
type PayloadConverter struct{}

// NewPayloadConverter creates a pointer to a new PayloadConverter which satisfies the PayloadConverter interface
func NewPayloadConverter() *PayloadConverter {
	return &PayloadConverter{}
}

// Convert method is used to convert a bitcoin BlockPayload to an IPLDPayload
// Satisfies the shared.PayloadConverter interface
func (pc *PayloadConverter) Convert(payload shared.RawChainData) (shared.StreamedIPLDs, error) {
	btcBlockPayload, ok := payload.(BlockPayload)
	if !ok {
		return nil, fmt.Errorf("btc converter: expected payload type %T got %T", BlockPayload{}, payload)
	}
	txMeta := make([]TxModel, len(btcBlockPayload.Txs))
	for _, tx := range btcBlockPayload.Txs {
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
		BlockPayload: btcBlockPayload,
		TxMetaData:   txMeta,
	}, nil
}
