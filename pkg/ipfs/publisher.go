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
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	rlp2 "github.com/ethereum/go-ethereum/rlp"
	"github.com/ipfs/go-ipfs/plugin/loader"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_state_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_storage_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

// IPLDPublisher is the interface for publishing an IPLD payload
type IPLDPublisher interface {
	Publish(payload *IPLDPayload) (*CIDPayload, error)
}

// Publisher is the underlying struct for the IPLDPublisher interface
type Publisher struct {
	Node              *ipfs.IPFS
	HeaderPutter      *eth_block_header.BlockHeaderDagPutter
	TransactionPutter *eth_block_transactions.BlockTransactionsDagPutter
	ReceiptPutter     *eth_block_receipts.EthBlockReceiptDagPutter
	StatePutter       *eth_state_trie.StateTrieDagPutter
	StoragePutter     *eth_storage_trie.StorageTrieDagPutter
}

// NewIPLDPublisher creates a pointer to a new Publisher which satisfies the IPLDPublisher interface
func NewIPLDPublisher(ipfsPath string) (*Publisher, error) {
	l, err := loader.NewPluginLoader("~/.ipfs/plugins")
	if err != nil {
		return nil, err
	}
	err = l.Initialize()
	if err != nil {
		return nil, err
	}
	err = l.Inject()
	if err != nil {
		return nil, err
	}
	node, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		Node:              node,
		HeaderPutter:      eth_block_header.NewBlockHeaderDagPutter(node, rlp.RlpDecoder{}),
		TransactionPutter: eth_block_transactions.NewBlockTransactionsDagPutter(node),
		ReceiptPutter:     eth_block_receipts.NewEthBlockReceiptDagPutter(node),
		StatePutter:       eth_state_trie.NewStateTrieDagPutter(node),
		StoragePutter:     eth_storage_trie.NewStorageTrieDagPutter(node),
	}, nil
}

// Publish publishes an IPLDPayload to IPFS and returns the corresponding CIDPayload
func (pub *Publisher) Publish(payload *IPLDPayload) (*CIDPayload, error) {
	// Process and publish headers
	headerCid, err := pub.publishHeaders(payload.HeaderRLP)
	if err != nil {
		return nil, err
	}

	// Process and publish uncles
	uncleCids := make(map[common.Hash]string)
	for _, uncle := range payload.BlockBody.Uncles {
		uncleRlp, err := rlp2.EncodeToBytes(uncle)
		if err != nil {
			return nil, err
		}
		cid, err := pub.publishHeaders(uncleRlp)
		if err != nil {
			return nil, err
		}
		uncleCids[uncle.Hash()] = cid
	}

	// Process and publish transactions
	transactionCids, err := pub.publishTransactions(payload.BlockBody, payload.TrxMetaData)
	if err != nil {
		return nil, err
	}

	// Process and publish receipts
	receiptsCids, err := pub.publishReceipts(payload.Receipts, payload.ReceiptMetaData)
	if err != nil {
		return nil, err
	}

	// Process and publish state leafs
	stateNodeCids, err := pub.publishStateNodes(payload.StateNodes)
	if err != nil {
		return nil, err
	}

	// Process and publish storage leafs
	storageNodeCids, err := pub.publishStorageNodes(payload.StorageNodes)
	if err != nil {
		return nil, err
	}

	// Package CIDs and their metadata into a single struct
	return &CIDPayload{
		BlockHash:       payload.BlockHash,
		BlockNumber:     payload.BlockNumber.String(),
		HeaderCID:       headerCid,
		UncleCIDS:       uncleCids,
		TransactionCIDs: transactionCids,
		ReceiptCIDs:     receiptsCids,
		StateNodeCIDs:   stateNodeCids,
		StorageNodeCIDs: storageNodeCids,
	}, nil
}

func (pub *Publisher) publishHeaders(headerRLP []byte) (string, error) {
	headerCids, err := pub.HeaderPutter.DagPut(headerRLP)
	if err != nil {
		return "", err
	}
	if len(headerCids) != 1 {
		return "", errors.New("single CID expected to be returned for header")
	}
	return headerCids[0], nil
}

func (pub *Publisher) publishTransactions(blockBody *types.Body, trxMeta []*TrxMetaData) (map[common.Hash]*TrxMetaData, error) {
	/*
		println("publishing transactions")
		for _, trx := range blockBody.Transactions {
			println("trx value:")
			println(trx.Value().Int64())
		}
	*/
	transactionCids, err := pub.TransactionPutter.DagPut(blockBody)
	if err != nil {
		return nil, err
	}
	if len(transactionCids) != len(blockBody.Transactions) {
		return nil, errors.New("expected one CID for each transaction")
	}
	mappedTrxCids := make(map[common.Hash]*TrxMetaData, len(transactionCids))
	for i, trx := range blockBody.Transactions {
		mappedTrxCids[trx.Hash()] = trxMeta[i]
		mappedTrxCids[trx.Hash()].CID = transactionCids[i]
	}
	return mappedTrxCids, nil
}

func (pub *Publisher) publishReceipts(receipts types.Receipts, receiptMeta []*ReceiptMetaData) (map[common.Hash]*ReceiptMetaData, error) {
	receiptsCids, err := pub.ReceiptPutter.DagPut(receipts)
	if err != nil {
		return nil, err
	}
	if len(receiptsCids) != len(receipts) {
		return nil, errors.New("expected one CID for each receipt")
	}
	mappedRctCids := make(map[common.Hash]*ReceiptMetaData, len(receiptsCids))
	for i, rct := range receipts {
		mappedRctCids[rct.TxHash] = receiptMeta[i]
		mappedRctCids[rct.TxHash].CID = receiptsCids[i]
	}
	return mappedRctCids, nil
}

func (pub *Publisher) publishStateNodes(stateNodes map[common.Hash]StateNode) (map[common.Hash]StateNodeCID, error) {
	stateNodeCids := make(map[common.Hash]StateNodeCID)
	for addr, node := range stateNodes {
		stateNodeCid, err := pub.StatePutter.DagPut(node.Value)
		if err != nil {
			return nil, err
		}
		if len(stateNodeCid) != 1 {
			return nil, errors.New("single CID expected to be returned for state leaf")
		}
		stateNodeCids[addr] = StateNodeCID{
			CID:  stateNodeCid[0],
			Leaf: node.Leaf,
		}
	}
	return stateNodeCids, nil
}

func (pub *Publisher) publishStorageNodes(storageNodes map[common.Hash][]StorageNode) (map[common.Hash][]StorageNodeCID, error) {
	storageLeafCids := make(map[common.Hash][]StorageNodeCID)
	for addr, storageTrie := range storageNodes {
		storageLeafCids[addr] = make([]StorageNodeCID, 0)
		for _, node := range storageTrie {
			storageNodeCid, err := pub.StoragePutter.DagPut(node.Value)
			if err != nil {
				return nil, err
			}
			if len(storageNodeCid) != 1 {
				return nil, errors.New("single CID expected to be returned for storage leaf")
			}
			storageLeafCids[addr] = append(storageLeafCids[addr], StorageNodeCID{
				Key:  node.Key.Hex(),
				CID:  storageNodeCid[0],
				Leaf: node.Leaf,
			})
		}
	}
	return storageLeafCids, nil
}
