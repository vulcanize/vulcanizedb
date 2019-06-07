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

package ipfs

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/vulcanize/vulcanizedb/pkg/config"
)

// ResponseScreener is the inteface used to screen eth data and package appropriate data into a response payload
type ResponseScreener interface {
	ScreenResponse(streamFilters *config.Subscription, payload IPLDPayload) (*ResponsePayload, error)
}

// Screener is the underlying struct for the ReponseScreener interface
type Screener struct{}

// NewResponseScreener creates a new Screener satisfyign the ReponseScreener interface
func NewResponseScreener() *Screener {
	return &Screener{}
}

// ScreenResponse is used to filter through eth data to extract and package requested data into a ResponsePayload
func (s *Screener) ScreenResponse(streamFilters *config.Subscription, payload IPLDPayload) (*ResponsePayload, error) {
	response := new(ResponsePayload)
	err := s.filterHeaders(streamFilters, response, payload)
	if err != nil {
		return nil, err
	}
	txHashes, err := s.filterTransactions(streamFilters, response, payload)
	if err != nil {
		return nil, err
	}
	err = s.filerReceipts(streamFilters, response, payload, txHashes)
	if err != nil {
		return nil, err
	}
	err = s.filterState(streamFilters, response, payload)
	if err != nil {
		return nil, err
	}
	err = s.filterStorage(streamFilters, response, payload)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Screener) filterHeaders(streamFilters *config.Subscription, response *ResponsePayload, payload IPLDPayload) error {
	if !streamFilters.HeaderFilter.Off && checkRange(streamFilters.StartingBlock, streamFilters.EndingBlock, payload.BlockNumber.Int64()) {
		response.HeadersRlp = append(response.HeadersRlp, payload.HeaderRLP)
		if !streamFilters.HeaderFilter.FinalOnly {
			for _, uncle := range payload.BlockBody.Uncles {
				uncleRlp, err := rlp.EncodeToBytes(uncle)
				if err != nil {
					return err
				}
				response.UnclesRlp = append(response.UnclesRlp, uncleRlp)
			}
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

func (s *Screener) filterTransactions(streamFilters *config.Subscription, response *ResponsePayload, payload IPLDPayload) ([]common.Hash, error) {
	trxHashes := make([]common.Hash, 0, len(payload.BlockBody.Transactions))
	if !streamFilters.TrxFilter.Off && checkRange(streamFilters.StartingBlock, streamFilters.EndingBlock, payload.BlockNumber.Int64()) {
		for i, trx := range payload.BlockBody.Transactions {
			if checkTransactions(streamFilters.TrxFilter.Src, streamFilters.TrxFilter.Dst, payload.TrxMetaData[i].Src, payload.TrxMetaData[i].Dst) {
				trxBuffer := new(bytes.Buffer)
				err := trx.EncodeRLP(trxBuffer)
				if err != nil {
					return nil, err
				}
				trxHashes = append(trxHashes, trx.Hash())
				response.TransactionsRlp = append(response.TransactionsRlp, trxBuffer.Bytes())
			}
		}
	}
	return trxHashes, nil
}

func checkTransactions(wantedSrc, wantedDst []string, actualSrc, actualDst string) bool {
	// If we aren't filtering for any addresses, every transaction is a go
	if len(wantedDst) == 0 && len(wantedSrc) == 0 {
		return true
	}
	for _, src := range wantedSrc {
		if src == actualSrc {
			return true
		}
	}
	for _, dst := range wantedDst {
		if dst == actualDst {
			return true
		}
	}
	return false
}

func (s *Screener) filerReceipts(streamFilters *config.Subscription, response *ResponsePayload, payload IPLDPayload, trxHashes []common.Hash) error {
	if !streamFilters.ReceiptFilter.Off && checkRange(streamFilters.StartingBlock, streamFilters.EndingBlock, payload.BlockNumber.Int64()) {
		for i, receipt := range payload.Receipts {
			if checkReceipts(receipt, streamFilters.ReceiptFilter.Topic0s, payload.ReceiptMetaData[i].Topic0s, trxHashes) {
				receiptForStorage := (*types.ReceiptForStorage)(receipt)
				receiptBuffer := new(bytes.Buffer)
				err := receiptForStorage.EncodeRLP(receiptBuffer)
				if err != nil {
					return err
				}
				response.ReceiptsRlp = append(response.ReceiptsRlp, receiptBuffer.Bytes())
			}
		}
	}
	return nil
}

func checkReceipts(rct *types.Receipt, wantedTopics, actualTopics []string, wantedTrxHashes []common.Hash) bool {
	// If we aren't filtering for any topics, all topics are a go
	if len(wantedTopics) == 0 {
		return true
	}
	for _, wantedTrxHash := range wantedTrxHashes {
		if bytes.Equal(wantedTrxHash.Bytes(), rct.TxHash.Bytes()) {
			return true
		}
	}
	for _, wantedTopic := range wantedTopics {
		for _, actualTopic := range actualTopics {
			if wantedTopic == actualTopic {
				return true
			}
		}
	}
	return false
}

func (s *Screener) filterState(streamFilters *config.Subscription, response *ResponsePayload, payload IPLDPayload) error {
	response.StateNodesRlp = make(map[common.Hash][]byte)
	if !streamFilters.StateFilter.Off && checkRange(streamFilters.StartingBlock, streamFilters.EndingBlock, payload.BlockNumber.Int64()) {
		keyFilters := make([]common.Hash, 0, len(streamFilters.StateFilter.Addresses))
		for _, addr := range streamFilters.StateFilter.Addresses {
			keyFilter := AddressToKey(common.HexToAddress(addr))
			keyFilters = append(keyFilters, keyFilter)
		}
		for key, stateNode := range payload.StateNodes {
			if checkNodeKeys(keyFilters, key) {
				if stateNode.Leaf || streamFilters.StateFilter.IntermediateNodes {
					response.StateNodesRlp[key] = stateNode.Value
				}
			}
		}
	}
	return nil
}

func checkNodeKeys(wantedKeys []common.Hash, actualKey common.Hash) bool {
	// If we aren't filtering for any specific keys, all nodes are a go
	if len(wantedKeys) == 0 {
		return true
	}
	for _, key := range wantedKeys {
		if bytes.Equal(key.Bytes(), actualKey.Bytes()) {
			return true
		}
	}
	return false
}

func (s *Screener) filterStorage(streamFilters *config.Subscription, response *ResponsePayload, payload IPLDPayload) error {
	if !streamFilters.StorageFilter.Off && checkRange(streamFilters.StartingBlock, streamFilters.EndingBlock, payload.BlockNumber.Int64()) {
		stateKeyFilters := make([]common.Hash, 0, len(streamFilters.StorageFilter.Addresses))
		for _, addr := range streamFilters.StorageFilter.Addresses {
			keyFilter := AddressToKey(common.HexToAddress(addr))
			stateKeyFilters = append(stateKeyFilters, keyFilter)
		}
		storageKeyFilters := make([]common.Hash, 0, len(streamFilters.StorageFilter.StorageKeys))
		for _, store := range streamFilters.StorageFilter.StorageKeys {
			keyFilter := HexToKey(store)
			storageKeyFilters = append(storageKeyFilters, keyFilter)
		}
		for stateKey, storageNodes := range payload.StorageNodes {
			if checkNodeKeys(stateKeyFilters, stateKey) {
				response.StorageNodesRlp[stateKey] = make(map[common.Hash][]byte)
				for _, storageNode := range storageNodes {
					if checkNodeKeys(storageKeyFilters, storageNode.Key) {
						response.StorageNodesRlp[stateKey][storageNode.Key] = storageNode.Value
					}
				}
			}
		}
	}
	return nil
}
