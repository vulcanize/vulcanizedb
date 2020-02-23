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

	"github.com/multiformats/go-multihash"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// ResponseFilterer satisfies the ResponseFilterer interface for bitcoin
type ResponseFilterer struct{}

// NewResponseFilterer creates a new Filterer satisfying the ResponseFilterer interface
func NewResponseFilterer() *ResponseFilterer {
	return &ResponseFilterer{}
}

// Filter is used to filter through btc data to extract and package requested data into a Payload
func (s *ResponseFilterer) Filter(filter shared.SubscriptionSettings, payload shared.ConvertedData) (shared.IPLDs, error) {
	btcFilters, ok := filter.(*SubscriptionSettings)
	if !ok {
		return IPLDs{}, fmt.Errorf("btc filterer expected filter type %T got %T", &SubscriptionSettings{}, filter)
	}
	btcPayload, ok := payload.(ConvertedPayload)
	if !ok {
		return IPLDs{}, fmt.Errorf("btc filterer expected payload type %T got %T", ConvertedPayload{}, payload)
	}
	height := int64(btcPayload.BlockPayload.BlockHeight)
	if checkRange(btcFilters.Start.Int64(), btcFilters.End.Int64(), height) {
		response := new(IPLDs)
		if err := s.filterHeaders(btcFilters.HeaderFilter, response, btcPayload); err != nil {
			return IPLDs{}, err
		}
		if err := s.filterTransactions(btcFilters.TxFilter, response, btcPayload); err != nil {
			return IPLDs{}, err
		}
		response.BlockNumber = big.NewInt(height)
		return *response, nil
	}
	return IPLDs{}, nil
}

func (s *ResponseFilterer) filterHeaders(headerFilter HeaderFilter, response *IPLDs, payload ConvertedPayload) error {
	if !headerFilter.Off {
		headerBuffer := new(bytes.Buffer)
		if err := payload.Header.Serialize(headerBuffer); err != nil {
			return err
		}
		data := headerBuffer.Bytes()
		cid, err := ipld.RawdataToCid(ipld.MBitcoinHeader, data, multihash.DBL_SHA2_256)
		if err != nil {
			return err
		}
		response.Header = ipfs.BlockModel{
			Data: data,
			CID:  cid.String(),
		}
	}
	return nil
}

func checkRange(start, end, actual int64) bool {
	if (end <= 0 || end >= actual) && start <= actual {
		return true
	}
	return false
}

func (s *ResponseFilterer) filterTransactions(trxFilter TxFilter, response *IPLDs, payload ConvertedPayload) error {
	if !trxFilter.Off {
		response.Transactions = make([]ipfs.BlockModel, 0, len(payload.TxMetaData))
		for i, txMeta := range payload.TxMetaData {
			if checkTransaction(txMeta, trxFilter) {
				trxBuffer := new(bytes.Buffer)
				if err := payload.Txs[i].MsgTx().Serialize(trxBuffer); err != nil {
					return err
				}
				data := trxBuffer.Bytes()
				cid, err := ipld.RawdataToCid(ipld.MBitcoinTx, data, multihash.DBL_SHA2_256)
				if err != nil {
					return err
				}
				response.Transactions = append(response.Transactions, ipfs.BlockModel{
					Data: data,
					CID:  cid.String(),
				})
			}
		}
	}
	return nil
}

// checkTransaction returns true if the provided transaction has a hit on the filter
func checkTransaction(txMeta TxModelWithInsAndOuts, txFilter TxFilter) bool {
	passesSegwitFilter := false
	if !txFilter.Segwit || (txFilter.Segwit && txMeta.SegWit) {
		passesSegwitFilter = true
	}
	passesMultiSigFilter := !txFilter.MultiSig
	if txFilter.MultiSig {
		for _, out := range txMeta.TxOutputs {
			if out.RequiredSigs > 1 {
				passesMultiSigFilter = true
			}
		}
	}
	passesWitnessFilter := len(txFilter.WitnessHashes) == 0
	for _, wantedWitnessHash := range txFilter.WitnessHashes {
		if wantedWitnessHash == txMeta.WitnessHash {
			passesWitnessFilter = true
		}
	}
	passesAddressFilter := len(txFilter.Addresses) == 0
	for _, wantedAddress := range txFilter.Addresses {
		for _, out := range txMeta.TxOutputs {
			for _, actualAddress := range out.Addresses {
				if wantedAddress == actualAddress {
					passesAddressFilter = true
				}
			}
		}
	}
	passesIndexFilter := len(txFilter.Indexes) == 0
	for _, wantedIndex := range txFilter.Indexes {
		if wantedIndex == txMeta.Index {
			passesIndexFilter = true
		}
	}
	passesPkScriptClassFilter := len(txFilter.PkScriptClasses) == 0
	for _, wantedPkScriptClass := range txFilter.PkScriptClasses {
		for _, out := range txMeta.TxOutputs {
			if out.ScriptClass == wantedPkScriptClass {
				passesPkScriptClassFilter = true
			}
		}
	}
	return passesSegwitFilter && passesMultiSigFilter && passesWitnessFilter && passesAddressFilter && passesIndexFilter && passesPkScriptClassFilter
}
