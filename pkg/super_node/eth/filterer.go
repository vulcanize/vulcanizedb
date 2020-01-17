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

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/config"
)

// ResponseFilterer satisfies the ResponseFilterer interface for ethereum
type ResponseFilterer struct{}

// NewResponseFilterer creates a new Filterer satisfying the ResponseFilterer interface
func NewResponseFilterer() *ResponseFilterer {
	return &ResponseFilterer{}
}

// Filter is used to filter through eth data to extract and package requested data into a Payload
func (s *ResponseFilterer) Filter(filter, payload interface{}) (interface{}, error) {
	ethFilters, ok := filter.(*config.EthSubscription)
	if !ok {
		return StreamPayload{}, fmt.Errorf("eth filterer expected filter type %T got %T", &config.EthSubscription{}, filter)
	}
	ethPayload, ok := payload.(*IPLDPayload)
	if !ok {
		return StreamPayload{}, fmt.Errorf("eth filterer expected payload type %T got %T", &IPLDPayload{}, payload)
	}
	if checkRange(ethFilters.Start.Int64(), ethFilters.End.Int64(), ethPayload.Block.Number().Int64()) {
		response := new(StreamPayload)
		if err := s.filterHeaders(ethFilters.HeaderFilter, response, ethPayload); err != nil {
			return StreamPayload{}, err
		}
		txHashes, err := s.filterTransactions(ethFilters.TxFilter, response, ethPayload)
		if err != nil {
			return StreamPayload{}, err
		}
		if err := s.filerReceipts(ethFilters.ReceiptFilter, response, ethPayload, txHashes); err != nil {
			return StreamPayload{}, err
		}
		if err := s.filterState(ethFilters.StateFilter, response, ethPayload); err != nil {
			return StreamPayload{}, err
		}
		if err := s.filterStorage(ethFilters.StorageFilter, response, ethPayload); err != nil {
			return StreamPayload{}, err
		}
		response.BlockNumber = ethPayload.Block.Number()
		return *response, nil
	}
	return StreamPayload{}, nil
}

func (s *ResponseFilterer) filterHeaders(headerFilter config.HeaderFilter, response *StreamPayload, payload *IPLDPayload) error {
	if !headerFilter.Off {
		response.HeadersRlp = append(response.HeadersRlp, payload.HeaderRLP)
		if headerFilter.Uncles {
			response.UnclesRlp = make([][]byte, 0, len(payload.Block.Body().Uncles))
			for _, uncle := range payload.Block.Body().Uncles {
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

func (s *ResponseFilterer) filterTransactions(trxFilter config.TxFilter, response *StreamPayload, payload *IPLDPayload) ([]common.Hash, error) {
	trxHashes := make([]common.Hash, 0, len(payload.Block.Body().Transactions))
	if !trxFilter.Off {
		for i, trx := range payload.Block.Body().Transactions {
			if checkTransactions(trxFilter.Src, trxFilter.Dst, payload.TrxMetaData[i].Src, payload.TrxMetaData[i].Dst) {
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

func (s *ResponseFilterer) filerReceipts(receiptFilter config.ReceiptFilter, response *StreamPayload, payload *IPLDPayload, trxHashes []common.Hash) error {
	if !receiptFilter.Off {
		for i, receipt := range payload.Receipts {
			if checkReceipts(receipt, receiptFilter.Topic0s, payload.ReceiptMetaData[i].Topic0s, receiptFilter.Contracts, payload.ReceiptMetaData[i].Contract, trxHashes, receiptFilter.MatchTxs) {
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

func checkReceipts(rct *types.Receipt, wantedTopics, actualTopics, wantedContracts []string, actualContract string, wantedTrxHashes []common.Hash, matchTxs bool) bool {
	// If we aren't filtering for any topics, contracts, or corresponding trxs then all receipts are a go
	if len(wantedTopics) == 0 && len(wantedContracts) == 0 && (len(wantedTrxHashes) == 0 || !matchTxs) {
		return true
	}
	// No matter what filters we have, we keep receipts for specific trxs we are interested in
	if matchTxs {
		for _, wantedTrxHash := range wantedTrxHashes {
			if bytes.Equal(wantedTrxHash.Bytes(), rct.TxHash.Bytes()) {
				return true
			}
		}
	}

	if len(wantedContracts) == 0 {
		// We keep all receipts that have logs we are interested in
		for _, wantedTopic := range wantedTopics {
			for _, actualTopic := range actualTopics {
				if wantedTopic == actualTopic {
					return true
				}
			}
		}
	} else { // We keep all receipts that belong to one of the specified contracts if we aren't filtering on topics
		for _, wantedContract := range wantedContracts {
			if wantedContract == actualContract {
				if len(wantedTopics) == 0 {
					return true
				}
				// Or if we have contracts and topics to filter on we only keep receipts that satisfy both conditions
				for _, wantedTopic := range wantedTopics {
					for _, actualTopic := range actualTopics {
						if wantedTopic == actualTopic {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (s *ResponseFilterer) filterState(stateFilter config.StateFilter, response *StreamPayload, payload *IPLDPayload) error {
	if !stateFilter.Off {
		response.StateNodesRlp = make(map[common.Hash][]byte)
		keyFilters := make([]common.Hash, 0, len(stateFilter.Addresses))
		for _, addr := range stateFilter.Addresses {
			keyFilter := AddressToKey(common.HexToAddress(addr))
			keyFilters = append(keyFilters, keyFilter)
		}
		for _, stateNode := range payload.StateNodes {
			if checkNodeKeys(keyFilters, stateNode.Key) {
				if stateNode.Leaf || stateFilter.IntermediateNodes {
					response.StateNodesRlp[stateNode.Key] = stateNode.Value
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

func (s *ResponseFilterer) filterStorage(storageFilter config.StorageFilter, response *StreamPayload, payload *IPLDPayload) error {
	if !storageFilter.Off {
		response.StorageNodesRlp = make(map[common.Hash]map[common.Hash][]byte)
		stateKeyFilters := make([]common.Hash, 0, len(storageFilter.Addresses))
		for _, addr := range storageFilter.Addresses {
			keyFilter := AddressToKey(common.HexToAddress(addr))
			stateKeyFilters = append(stateKeyFilters, keyFilter)
		}
		storageKeyFilters := make([]common.Hash, 0, len(storageFilter.StorageKeys))
		for _, store := range storageFilter.StorageKeys {
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
