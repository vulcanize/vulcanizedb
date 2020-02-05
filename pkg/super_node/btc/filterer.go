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
	"math/big"

	"github.com/btcsuite/btcutil"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// ResponseFilterer satisfies the ResponseFilterer interface for ethereum
type ResponseFilterer struct{}

// NewResponseFilterer creates a new Filterer satisfying the ResponseFilterer interface
func NewResponseFilterer() *ResponseFilterer {
	return &ResponseFilterer{}
}

// Filter is used to filter through eth data to extract and package requested data into a Payload
func (s *ResponseFilterer) Filter(filter shared.SubscriptionSettings, payload shared.StreamedIPLDs) (shared.ServerResponse, error) {
	btcFilters, ok := filter.(*SubscriptionSettings)
	if !ok {
		return StreamResponse{}, fmt.Errorf("eth filterer expected filter type %T got %T", &SubscriptionSettings{}, filter)
	}
	btcPayload, ok := payload.(IPLDPayload)
	if !ok {
		return StreamResponse{}, fmt.Errorf("eth filterer expected payload type %T got %T", IPLDPayload{}, payload)
	}
	height := int64(btcPayload.Height)
	if checkRange(btcFilters.Start.Int64(), btcFilters.End.Int64(), height) {
		response := new(StreamResponse)
		if err := s.filterHeaders(btcFilters.HeaderFilter, response, btcPayload); err != nil {
			return StreamResponse{}, err
		}
		if err := s.filterTransactions(btcFilters.TxFilter, response, btcPayload); err != nil {
			return StreamResponse{}, err
		}
		response.BlockNumber = big.NewInt(height)
		return *response, nil
	}
	return StreamResponse{}, nil
}

func (s *ResponseFilterer) filterHeaders(headerFilter HeaderFilter, response *StreamResponse, payload IPLDPayload) error {
	if !headerFilter.Off {
		headerBuffer := new(bytes.Buffer)
		if err := payload.Header.Serialize(headerBuffer); err != nil {
			return err
		}
		response.SerializedHeaders = append(response.SerializedHeaders, headerBuffer.Bytes())
	}
	return nil
}

func checkRange(start, end, actual int64) bool {
	if (end <= 0 || end >= actual) && start <= actual {
		return true
	}
	return false
}

func (s *ResponseFilterer) filterTransactions(trxFilter TxFilter, response *StreamResponse, payload IPLDPayload) error {
	if !trxFilter.Off {
		for _, trx := range payload.Txs {
			if checkTransaction(trx, trxFilter) {
				trxBuffer := new(bytes.Buffer)
				if err := trx.MsgTx().Serialize(trxBuffer); err != nil {
					return err
				}
				response.SerializedTxs = append(response.SerializedTxs, trxBuffer.Bytes())
			}
		}
	}
	return nil
}

// checkTransaction returns true if the provided transaction has a hit on the filter
func checkTransaction(trx *btcutil.Tx, txFilter TxFilter) bool {
	panic("implement me")
}
