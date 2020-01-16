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

package super_node

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// ResponseFilterer is the inteface used to screen eth data and package appropriate data into a response payload
type ResponseFilterer interface {
	FilterResponse(streamFilters config.Subscription, payload ipfs.IPLDPayload) (streamer.SuperNodePayload, error)
}

// Filterer is the underlying struct for the ResponseFilterer interface
type Filterer struct{}

// NewResponseFilterer creates a new Filterer satisfying the ResponseFilterer interface
func NewResponseFilterer() *Filterer {
	return &Filterer{}
}

// FilterResponse is used to filter through eth data to extract and package requested data into a Payload
func (s *Filterer) FilterResponse(streamFilters config.Subscription, payload ipfs.IPLDPayload) (streamer.SuperNodePayload, error) {
	if checkRange(streamFilters.StartingBlock.Int64(), streamFilters.EndingBlock.Int64(), payload.BlockNumber.Int64()) {
		response := new(streamer.SuperNodePayload)
		if err := s.filterHeaders(streamFilters.HeaderFilter, response, payload); err != nil {
			return streamer.SuperNodePayload{}, err
		}
		txHashes, err := s.filterTransactions(streamFilters.TrxFilter, response, payload)
		if err != nil {
			return streamer.SuperNodePayload{}, err
		}
		if err := s.filerReceipts(streamFilters.ReceiptFilter, response, payload, txHashes); err != nil {
			return streamer.SuperNodePayload{}, err
		}
		if err := s.filterState(streamFilters.StateFilter, response, payload); err != nil {
			return streamer.SuperNodePayload{}, err
		}
		if err := s.filterStorage(streamFilters.StorageFilter, response, payload); err != nil {
			return streamer.SuperNodePayload{}, err
		}
		response.BlockNumber = payload.BlockNumber
		return *response, nil
	}
	return streamer.SuperNodePayload{}, nil
}

func (s *Filterer) filterHeaders(headerFilter config.HeaderFilter, response *streamer.SuperNodePayload, payload ipfs.IPLDPayload) error {
	if !headerFilter.Off {
		response.HeadersRlp = append(response.HeadersRlp, payload.HeaderRLP)
		if headerFilter.Uncles {
			response.UnclesRlp = make([][]byte, 0, len(payload.BlockBody.Uncles))
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

func (s *Filterer) filterTransactions(trxFilter config.TrxFilter, response *streamer.SuperNodePayload, payload ipfs.IPLDPayload) ([]common.Hash, error) {
	trxHashes := make([]common.Hash, 0, len(payload.BlockBody.Transactions))
	if !trxFilter.Off {
		for i, trx := range payload.BlockBody.Transactions {
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

func (s *Filterer) filerReceipts(receiptFilter config.ReceiptFilter, response *streamer.SuperNodePayload, payload ipfs.IPLDPayload, trxHashes []common.Hash) error {
	if !receiptFilter.Off {
		for i, receipt := range payload.Receipts {
			if checkReceipts(receipt, receiptFilter.Topic0s, payload.ReceiptMetaData[i].Topic0s, receiptFilter.Contracts, payload.ReceiptMetaData[i].ContractAddress, trxHashes, receiptFilter.MatchTxs) {
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

func (s *Filterer) filterState(stateFilter config.StateFilter, response *streamer.SuperNodePayload, payload ipfs.IPLDPayload) error {
	if !stateFilter.Off {
		response.StateNodesRlp = make(map[common.Hash][]byte)
		keyFilters := make([]common.Hash, 0, len(stateFilter.Addresses))
		for _, addr := range stateFilter.Addresses {
			keyFilter := ipfs.AddressToKey(common.HexToAddress(addr))
			keyFilters = append(keyFilters, keyFilter)
		}
		for key, stateNode := range payload.StateNodes {
			if checkNodeKeys(keyFilters, key) {
				if stateNode.Leaf || stateFilter.IntermediateNodes {
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

func (s *Filterer) filterStorage(storageFilter config.StorageFilter, response *streamer.SuperNodePayload, payload ipfs.IPLDPayload) error {
	if !storageFilter.Off {
		response.StorageNodesRlp = make(map[common.Hash]map[common.Hash][]byte)
		stateKeyFilters := make([]common.Hash, 0, len(storageFilter.Addresses))
		for _, addr := range storageFilter.Addresses {
			keyFilter := ipfs.AddressToKey(common.HexToAddress(addr))
			stateKeyFilters = append(stateKeyFilters, keyFilter)
		}
		storageKeyFilters := make([]common.Hash, 0, len(storageFilter.StorageKeys))
		for _, store := range storageFilter.StorageKeys {
			keyFilter := ipfs.HexToKey(store)
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
