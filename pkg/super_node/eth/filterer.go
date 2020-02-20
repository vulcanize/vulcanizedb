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
	"github.com/multiformats/go-multihash"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// ResponseFilterer satisfies the ResponseFilterer interface for ethereum
type ResponseFilterer struct{}

// NewResponseFilterer creates a new Filterer satisfying the ResponseFilterer interface
func NewResponseFilterer() *ResponseFilterer {
	return &ResponseFilterer{}
}

// Filter is used to filter through eth data to extract and package requested data into a Payload
func (s *ResponseFilterer) Filter(filter shared.SubscriptionSettings, payload shared.ConvertedData) (shared.IPLDs, error) {
	ethFilters, ok := filter.(*SubscriptionSettings)
	if !ok {
		return IPLDs{}, fmt.Errorf("eth filterer expected filter type %T got %T", &SubscriptionSettings{}, filter)
	}
	ethPayload, ok := payload.(ConvertedPayload)
	if !ok {
		return IPLDs{}, fmt.Errorf("eth filterer expected payload type %T got %T", ConvertedPayload{}, payload)
	}
	if checkRange(ethFilters.Start.Int64(), ethFilters.End.Int64(), ethPayload.Block.Number().Int64()) {
		response := new(IPLDs)
		if err := s.filterHeaders(ethFilters.HeaderFilter, response, ethPayload); err != nil {
			return IPLDs{}, err
		}
		txHashes, err := s.filterTransactions(ethFilters.TxFilter, response, ethPayload)
		if err != nil {
			return IPLDs{}, err
		}
		var filterTxs []common.Hash
		if ethFilters.ReceiptFilter.MatchTxs {
			filterTxs = txHashes
		}
		if err := s.filerReceipts(ethFilters.ReceiptFilter, response, ethPayload, filterTxs); err != nil {
			return IPLDs{}, err
		}
		if err := s.filterState(ethFilters.StateFilter, response, ethPayload); err != nil {
			return IPLDs{}, err
		}
		if err := s.filterStorage(ethFilters.StorageFilter, response, ethPayload); err != nil {
			return IPLDs{}, err
		}
		response.BlockNumber = ethPayload.Block.Number()
		return *response, nil
	}
	return IPLDs{}, nil
}

func (s *ResponseFilterer) filterHeaders(headerFilter HeaderFilter, response *IPLDs, payload ConvertedPayload) error {
	if !headerFilter.Off {
		headerRLP, err := rlp.EncodeToBytes(payload.Block.Header())
		if err != nil {
			return err
		}
		cid, err := ipld.RawdataToCid(ipld.MEthHeader, headerRLP, multihash.KECCAK_256)
		if err != nil {
			return err
		}
		response.Headers = append(response.Headers, ipfs.BlockModel{
			Data: headerRLP,
			CID:  cid.String(),
		})
		if headerFilter.Uncles {
			response.Uncles = make([]ipfs.BlockModel, len(payload.Block.Body().Uncles))
			for i, uncle := range payload.Block.Body().Uncles {
				uncleRlp, err := rlp.EncodeToBytes(uncle)
				if err != nil {
					return err
				}
				cid, err := ipld.RawdataToCid(ipld.MEthHeader, uncleRlp, multihash.KECCAK_256)
				if err != nil {
					return err
				}
				response.Uncles[i] = ipfs.BlockModel{
					Data: uncleRlp,
					CID:  cid.String(),
				}
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

func (s *ResponseFilterer) filterTransactions(trxFilter TxFilter, response *IPLDs, payload ConvertedPayload) ([]common.Hash, error) {
	var trxHashes []common.Hash
	if !trxFilter.Off {
		trxLen := len(payload.Block.Body().Transactions)
		trxHashes = make([]common.Hash, 0, trxLen)
		response.Transactions = make([]ipfs.BlockModel, 0, trxLen)
		for i, trx := range payload.Block.Body().Transactions {
			if checkTransactionAddrs(trxFilter.Src, trxFilter.Dst, payload.TxMetaData[i].Src, payload.TxMetaData[i].Dst) {
				trxBuffer := new(bytes.Buffer)
				if err := trx.EncodeRLP(trxBuffer); err != nil {
					return nil, err
				}
				data := trxBuffer.Bytes()
				cid, err := ipld.RawdataToCid(ipld.MEthTx, data, multihash.KECCAK_256)
				if err != nil {
					return nil, err
				}
				response.Transactions = append(response.Transactions, ipfs.BlockModel{
					Data: data,
					CID:  cid.String(),
				})
				trxHashes = append(trxHashes, trx.Hash())
			}
		}
	}
	return trxHashes, nil
}

// checkTransactionAddrs returns true if either the transaction src and dst are one of the wanted src and dst addresses
func checkTransactionAddrs(wantedSrc, wantedDst []string, actualSrc, actualDst string) bool {
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

func (s *ResponseFilterer) filerReceipts(receiptFilter ReceiptFilter, response *IPLDs, payload ConvertedPayload, trxHashes []common.Hash) error {
	if !receiptFilter.Off {
		response.Receipts = make([]ipfs.BlockModel, 0, len(payload.Receipts))
		for i, receipt := range payload.Receipts {
			// topics is always length 4
			topics := [][]string{payload.ReceiptMetaData[i].Topic0s, payload.ReceiptMetaData[i].Topic1s, payload.ReceiptMetaData[i].Topic2s, payload.ReceiptMetaData[i].Topic3s}
			if checkReceipts(receipt, receiptFilter.Topics, topics, receiptFilter.Contracts, payload.ReceiptMetaData[i].Contract, trxHashes) {
				receiptBuffer := new(bytes.Buffer)
				if err := receipt.EncodeRLP(receiptBuffer); err != nil {
					return err
				}
				data := receiptBuffer.Bytes()
				cid, err := ipld.RawdataToCid(ipld.MEthTxReceipt, data, multihash.KECCAK_256)
				if err != nil {
					return err
				}
				response.Receipts = append(response.Receipts, ipfs.BlockModel{
					Data: data,
					CID:  cid.String(),
				})
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

// filterMatch returns true if the actualTopics conform to the wantedTopics filter
func filterMatch(wantedTopics, actualTopics [][]string) bool {
	// actualTopics should always be length 4, but the members can be nil slices
	matches := 0
	for i, actualTopicSet := range actualTopics {
		if i < len(wantedTopics) && len(wantedTopics[i]) > 0 {
			// If we have topics in this filter slot, count as a match if one of the topics matches
			matches += slicesShareString(actualTopicSet, wantedTopics[i])
		} else {
			// Filter slot is either empty or doesn't exist => not matching any topics at this slot => counts as a match
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

func (s *ResponseFilterer) filterState(stateFilter StateFilter, response *IPLDs, payload ConvertedPayload) error {
	if !stateFilter.Off {
		response.StateNodes = make([]StateNode, 0, len(payload.StateNodes))
		keyFilters := make([]common.Hash, len(stateFilter.Addresses))
		for i, addr := range stateFilter.Addresses {
			keyFilters[i] = crypto.Keccak256Hash(common.HexToAddress(addr).Bytes())
		}
		for _, stateNode := range payload.StateNodes {
			if checkNodeKeys(keyFilters, stateNode.Key) {
				if stateNode.Leaf || stateFilter.IntermediateNodes {
					cid, err := ipld.RawdataToCid(ipld.MEthStateTrie, stateNode.Value, multihash.KECCAK_256)
					if err != nil {
						return err
					}
					response.StateNodes = append(response.StateNodes, StateNode{
						StateTrieKey: stateNode.Key,
						IPLD: ipfs.BlockModel{
							Data: stateNode.Value,
							CID:  cid.String(),
						},
						Leaf: stateNode.Leaf,
					})
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

func (s *ResponseFilterer) filterStorage(storageFilter StorageFilter, response *IPLDs, payload ConvertedPayload) error {
	if !storageFilter.Off {
		response.StorageNodes = make([]StorageNode, 0)
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
				for _, storageNode := range storageNodes {
					if checkNodeKeys(storageKeyFilters, storageNode.Key) {
						cid, err := ipld.RawdataToCid(ipld.MEthStorageTrie, storageNode.Value, multihash.KECCAK_256)
						if err != nil {
							return err
						}
						response.StorageNodes = append(response.StorageNodes, StorageNode{
							StateTrieKey:   stateKey,
							StorageTrieKey: storageNode.Key,
							IPLD: ipfs.BlockModel{
								Data: storageNode.Value,
								CID:  cid.String(),
							},
							Leaf: storageNode.Leaf,
						})
					}
				}
			}
		}
	}
	return nil
}
