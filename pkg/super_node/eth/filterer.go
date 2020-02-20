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
	"github.com/ethereum/go-ethereum/crypto"
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
		var filterTxs []common.Hash
		if ethFilters.ReceiptFilter.MatchTxs {
			filterTxs = txHashes
		}
		if err := s.filerReceipts(ethFilters.ReceiptFilter, response, ethPayload, filterTxs); err != nil {
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
				if err := trx.EncodeRLP(trxBuffer); err != nil {
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
			// topics is always length 4
			topics := [][]string{payload.ReceiptMetaData[i].Topic0s, payload.ReceiptMetaData[i].Topic1s, payload.ReceiptMetaData[i].Topic2s, payload.ReceiptMetaData[i].Topic3s}
			if checkReceipts(receipt, receiptFilter.Topics, topics, receiptFilter.Contracts, payload.ReceiptMetaData[i].Contract, trxHashes) {
				receiptForStorage := (*types.ReceiptForStorage)(receipt)
				receiptBuffer := new(bytes.Buffer)
				if err := receiptForStorage.EncodeRLP(receiptBuffer); err != nil {
					return err
				}
				response.ReceiptsRlp = append(response.ReceiptsRlp, receiptBuffer.Bytes())
			}
		}
	}
	return nil
}

func checkReceipts(rct *types.Receipt, wantedTopics, actualTopics [][]string, wantedContracts []string, actualContract string, wantedTrxHashes []common.Hash) bool {
	// If we aren't filtering for any topics, contracts, or corresponding trxs then all receipts are a go
	if len(wantedTopics) == 0 && len(wantedContracts) == 0 && len(wantedTrxHashes) == 0 {
		return true
	}
	// Keep receipts that are from watched txs
	for _, wantedTrxHash := range wantedTrxHashes {
		if bytes.Equal(wantedTrxHash.Bytes(), rct.TxHash.Bytes()) {
			return true
		}
	}
	// If there are no wanted contract addresses, we keep all receipts that match the topic filter
	if len(wantedContracts) == 0 {
		if match := filterMatch(wantedTopics, actualTopics); match == true {
			return true
		}
	}
	// If there are wanted contract addresses to filter on
	for _, wantedAddr := range wantedContracts {
		// and this is an address of interest
		if wantedAddr == actualContract {
			// we keep the receipt if it matches on the topic filter
			if match := filterMatch(wantedTopics, actualTopics); match == true {
				return true
			}
		}
	}
	return false
}

func filterMatch(wantedTopics, actualTopics [][]string) bool {
	// actualTopics should always be length 4, members could be nil slices though
	lenWantedTopics := len(wantedTopics)
	matches := 0
	for i, actualTopicSet := range actualTopics {
		if i < lenWantedTopics {
			// If we have topics in this filter slot, count as a match if one of the topics matches
			if len(wantedTopics[i]) > 0 {
				matches += slicesShareString(actualTopicSet, wantedTopics[i])
			} else {
				// Filter slot is empty, not matching any topics at this slot => counts as a match
				matches++
			}
		} else {
			// Filter slot doesn't exist, not matching any topics at this slot => count as a match
			matches++
		}
	}
	if matches == 4 {
		return true
	}
	return false
}

// returns 1 if the two slices have a string in common, 0 if they do not
func slicesShareString(slice1, slice2 []string) int {
	for _, str1 := range slice1 {
		for _, str2 := range slice2 {
			if str1 == str2 {
				return 1
			}
		}
	}
	return 0
}

func (s *ResponseFilterer) filterState(stateFilter config.StateFilter, response *StreamPayload, payload *IPLDPayload) error {
	if !stateFilter.Off {
		response.StateNodesRlp = make(map[common.Hash][]byte)
		keyFilters := make([]common.Hash, len(stateFilter.Addresses))
		for i, addr := range stateFilter.Addresses {
			keyFilters[i] = crypto.Keccak256Hash(common.HexToAddress(addr).Bytes())
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
		stateKeyFilters := make([]common.Hash, len(storageFilter.Addresses))
		for i, addr := range storageFilter.Addresses {
			stateKeyFilters[i] = crypto.Keccak256Hash(common.HexToAddress(addr).Bytes())
		}
		storageKeyFilters := make([]common.Hash, len(storageFilter.StorageKeys))
		for i, store := range storageFilter.StorageKeys {
			storageKeyFilters[i] = common.HexToHash(store)
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
