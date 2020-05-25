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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"

	common2 "github.com/vulcanize/vulcanizedb/pkg/eth/converters/common"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/dag_putters"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// IPLDPublisher satisfies the IPLDPublisher for ethereum
type IPLDPublisher struct {
	HeaderPutter          ipfs.DagPutter
	TransactionPutter     ipfs.DagPutter
	TransactionTriePutter ipfs.DagPutter
	ReceiptPutter         ipfs.DagPutter
	ReceiptTriePutter     ipfs.DagPutter
	StatePutter           ipfs.DagPutter
	StoragePutter         ipfs.DagPutter
}

// NewIPLDPublisher creates a pointer to a new IPLDPublisher which satisfies the IPLDPublisher interface
func NewIPLDPublisher(ipfsPath string) (*IPLDPublisher, error) {
	node, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &IPLDPublisher{
		HeaderPutter:          dag_putters.NewEthBlockHeaderDagPutter(node),
		TransactionPutter:     dag_putters.NewEthTxsDagPutter(node),
		TransactionTriePutter: dag_putters.NewEthTxTrieDagPutter(node),
		ReceiptPutter:         dag_putters.NewEthReceiptDagPutter(node),
		ReceiptTriePutter:     dag_putters.NewEthRctTrieDagPutter(node),
		StatePutter:           dag_putters.NewEthStateDagPutter(node),
		StoragePutter:         dag_putters.NewEthStorageDagPutter(node),
	}, nil
}

// Publish publishes an IPLDPayload to IPFS and returns the corresponding CIDPayload
func (pub *IPLDPublisher) Publish(payload shared.ConvertedData) (shared.CIDsForIndexing, error) {
	ipldPayload, ok := payload.(ConvertedPayload)
	if !ok {
		return nil, fmt.Errorf("eth publisher expected payload type %T got %T", ConvertedPayload{}, payload)
	}
	// Generate the nodes for publishing
	headerNode, uncleNodes, txNodes, txTrieNodes, rctNodes, rctTrieNodes, err := ipld.FromBlockAndReceipts(ipldPayload.Block, ipldPayload.Receipts)
	if err != nil {
		return nil, err
	}

	// Process and publish headers
	headerCid, err := pub.publishHeader(headerNode)
	if err != nil {
		return nil, err
	}
	reward := common2.CalcEthBlockReward(ipldPayload.Block.Header(), ipldPayload.Block.Uncles(), ipldPayload.Block.Transactions(), ipldPayload.Receipts)
	header := HeaderModel{
		CID:             headerCid,
		ParentHash:      ipldPayload.Block.ParentHash().String(),
		BlockNumber:     ipldPayload.Block.Number().String(),
		BlockHash:       ipldPayload.Block.Hash().String(),
		TotalDifficulty: ipldPayload.TotalDifficulty.String(),
		Reward:          reward.String(),
		Bloom:           ipldPayload.Block.Bloom().Bytes(),
		StateRoot:       ipldPayload.Block.Root().String(),
		RctRoot:         ipldPayload.Block.ReceiptHash().String(),
		TxRoot:          ipldPayload.Block.TxHash().String(),
		UncleRoot:       ipldPayload.Block.UncleHash().String(),
		Timestamp:       ipldPayload.Block.Time(),
	}

	// Process and publish uncles
	uncleCids := make([]UncleModel, len(uncleNodes))
	for i, uncle := range uncleNodes {
		uncleCid, err := pub.publishHeader(uncle)
		if err != nil {
			return nil, err
		}
		uncleReward := common2.CalcUncleMinerReward(ipldPayload.Block.Number().Int64(), uncle.Number.Int64())
		uncleCids[i] = UncleModel{
			CID:        uncleCid,
			ParentHash: uncle.ParentHash.String(),
			BlockHash:  uncle.Hash().String(),
			Reward:     uncleReward.String(),
		}
	}

	// Process and publish transactions
	transactionCids, err := pub.publishTransactions(txNodes, txTrieNodes, ipldPayload.TxMetaData)
	if err != nil {
		return nil, err
	}

	// Process and publish receipts
	receiptsCids, err := pub.publishReceipts(rctNodes, rctTrieNodes, ipldPayload.ReceiptMetaData)
	if err != nil {
		return nil, err
	}

	// Process and publish state leafs
	stateNodeCids, stateAccounts, err := pub.publishStateNodes(ipldPayload.StateNodes)
	if err != nil {
		return nil, err
	}

	// Process and publish storage leafs
	storageNodeCids, err := pub.publishStorageNodes(ipldPayload.StorageNodes)
	if err != nil {
		return nil, err
	}

	// Package CIDs and their metadata into a single struct
	return &CIDPayload{
		HeaderCID:       header,
		UncleCIDs:       uncleCids,
		TransactionCIDs: transactionCids,
		ReceiptCIDs:     receiptsCids,
		StateNodeCIDs:   stateNodeCids,
		StorageNodeCIDs: storageNodeCids,
		StateAccounts:   stateAccounts,
	}, nil
}

func (pub *IPLDPublisher) generateBlockNodes(body *types.Block, receipts types.Receipts) (*ipld.EthHeader,
	[]*ipld.EthHeader, []*ipld.EthTx, []*ipld.EthTxTrie, []*ipld.EthReceipt, []*ipld.EthRctTrie, error) {
	return ipld.FromBlockAndReceipts(body, receipts)
}

func (pub *IPLDPublisher) publishHeader(header *ipld.EthHeader) (string, error) {
	return pub.HeaderPutter.DagPut(header)
}

func (pub *IPLDPublisher) publishTransactions(transactions []*ipld.EthTx, txTrie []*ipld.EthTxTrie, trxMeta []TxModel) ([]TxModel, error) {
	trxCids := make([]TxModel, len(transactions))
	for i, tx := range transactions {
		cid, err := pub.TransactionPutter.DagPut(tx)
		if err != nil {
			return nil, err
		}
		trxCids[i] = TxModel{
			CID:    cid,
			Index:  trxMeta[i].Index,
			TxHash: trxMeta[i].TxHash,
			Src:    trxMeta[i].Src,
			Dst:    trxMeta[i].Dst,
		}
	}
	for _, txNode := range txTrie {
		// We don't do anything with the tx trie cids atm
		if _, err := pub.TransactionTriePutter.DagPut(txNode); err != nil {
			return nil, err
		}
	}
	return trxCids, nil
}

func (pub *IPLDPublisher) publishReceipts(receipts []*ipld.EthReceipt, receiptTrie []*ipld.EthRctTrie, receiptMeta []ReceiptModel) (map[common.Hash]ReceiptModel, error) {
	rctCids := make(map[common.Hash]ReceiptModel)
	for i, rct := range receipts {
		cid, err := pub.ReceiptPutter.DagPut(rct)
		if err != nil {
			return nil, err
		}
		rctCids[rct.TxHash] = ReceiptModel{
			CID:          cid,
			Contract:     receiptMeta[i].Contract,
			ContractHash: receiptMeta[i].ContractHash,
			Topic0s:      receiptMeta[i].Topic0s,
			Topic1s:      receiptMeta[i].Topic1s,
			Topic2s:      receiptMeta[i].Topic2s,
			Topic3s:      receiptMeta[i].Topic3s,
			LogContracts: receiptMeta[i].LogContracts,
		}
	}
	for _, rctNode := range receiptTrie {
		// We don't do anything with the rct trie cids atm
		if _, err := pub.ReceiptTriePutter.DagPut(rctNode); err != nil {
			return nil, err
		}
	}
	return rctCids, nil
}

func (pub *IPLDPublisher) publishStateNodes(stateNodes []TrieNode) ([]StateNodeModel, map[string]StateAccountModel, error) {
	stateNodeCids := make([]StateNodeModel, 0, len(stateNodes))
	stateAccounts := make(map[string]StateAccountModel)
	for _, stateNode := range stateNodes {
		node, err := ipld.FromStateTrieRLP(stateNode.Value)
		if err != nil {
			return nil, nil, err
		}
		cid, err := pub.StatePutter.DagPut(node)
		if err != nil {
			return nil, nil, err
		}
		stateNodeCids = append(stateNodeCids, StateNodeModel{
			Path:     stateNode.Path,
			StateKey: stateNode.LeafKey.String(),
			CID:      cid,
			NodeType: ResolveFromNodeType(stateNode.Type),
		})
		// If we have a leaf, decode the account to extract additional metadata for indexing
		if stateNode.Type == statediff.Leaf {
			var i []interface{}
			if err := rlp.DecodeBytes(stateNode.Value, &i); err != nil {
				return nil, nil, err
			}
			if len(i) != 2 {
				return nil, nil, fmt.Errorf("IPLDPublisher expected state leaf node rlp to decode into two elements")
			}
			var account state.Account
			if err := rlp.DecodeBytes(i[1].([]byte), &account); err != nil {
				return nil, nil, err
			}
			// Map state account to the state path hash
			statePath := common.Bytes2Hex(stateNode.Path)
			stateAccounts[statePath] = StateAccountModel{
				Balance:     account.Balance.String(),
				Nonce:       account.Nonce,
				CodeHash:    account.CodeHash,
				StorageRoot: account.Root.String(),
			}
		}
	}
	return stateNodeCids, stateAccounts, nil
}

func (pub *IPLDPublisher) publishStorageNodes(storageNodes map[string][]TrieNode) (map[string][]StorageNodeModel, error) {
	storageLeafCids := make(map[string][]StorageNodeModel)
	for path, storageTrie := range storageNodes {
		storageLeafCids[path] = make([]StorageNodeModel, 0, len(storageTrie))
		for _, storageNode := range storageTrie {
			node, err := ipld.FromStorageTrieRLP(storageNode.Value)
			if err != nil {
				return nil, err
			}
			cid, err := pub.StoragePutter.DagPut(node)
			if err != nil {
				return nil, err
			}
			// Map storage node cids to the state path hash
			storageLeafCids[path] = append(storageLeafCids[path], StorageNodeModel{
				Path:       storageNode.Path,
				StorageKey: storageNode.LeafKey.Hex(),
				CID:        cid,
				NodeType:   ResolveFromNodeType(storageNode.Type),
			})
		}
	}
	return storageLeafCids, nil
}
