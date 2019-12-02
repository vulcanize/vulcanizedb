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
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_state_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_storage_trie"
	rlp2 "github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

// IPLDPublisher is the interface for publishing an IPLD payload
type IPLDPublisher interface {
	Publish(payload *IPLDPayload) (*CIDPayload, error)
}

// Publisher is the underlying struct for the IPLDPublisher interface
type Publisher struct {
	HeaderPutter      ipfs.DagPutter
	TransactionPutter ipfs.DagPutter
	ReceiptPutter     ipfs.DagPutter
	StatePutter       ipfs.DagPutter
	StoragePutter     ipfs.DagPutter
}

// NewIPLDPublisher creates a pointer to a new Publisher which satisfies the IPLDPublisher interface
func NewIPLDPublisher(ipfsPath string) (*Publisher, error) {
	node, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		HeaderPutter:      eth_block_header.NewBlockHeaderDagPutter(node, rlp2.RlpDecoder{}),
		TransactionPutter: eth_block_transactions.NewBlockTransactionsDagPutter(node),
		ReceiptPutter:     eth_block_receipts.NewEthBlockReceiptDagPutter(node),
		StatePutter:       eth_state_trie.NewStateTrieDagPutter(node),
		StoragePutter:     eth_storage_trie.NewStorageTrieDagPutter(node),
	}, nil
}

// Publish publishes an IPLDPayload to IPFS and returns the corresponding CIDPayload
func (pub *Publisher) Publish(payload *IPLDPayload) (*CIDPayload, error) {
	// Process and publish headers
	headerCid, headersErr := pub.publishHeaders(payload.HeaderRLP)
	if headersErr != nil {
		return nil, headersErr
	}

	// Process and publish uncles
	uncleCids := make(map[common.Hash]string)
	for _, uncle := range payload.BlockBody.Uncles {
		uncleRlp, encodeErr := rlp.EncodeToBytes(uncle)
		if encodeErr != nil {
			return nil, encodeErr
		}
		cid, unclesErr := pub.publishHeaders(uncleRlp)
		if unclesErr != nil {
			return nil, unclesErr
		}
		uncleCids[uncle.Hash()] = cid
	}

	// Process and publish transactions
	transactionCids, trxsErr := pub.publishTransactions(payload.BlockBody, payload.TrxMetaData)
	if trxsErr != nil {
		return nil, trxsErr
	}

	// Process and publish receipts
	receiptsCids, rctsErr := pub.publishReceipts(payload.Receipts, payload.ReceiptMetaData)
	if rctsErr != nil {
		return nil, rctsErr
	}

	// Process and publish state leafs
	stateNodeCids, stateErr := pub.publishStateNodes(payload.StateNodes)
	if stateErr != nil {
		return nil, stateErr
	}

	// Process and publish storage leafs
	storageNodeCids, storageErr := pub.publishStorageNodes(payload.StorageNodes)
	if storageErr != nil {
		return nil, storageErr
	}

	// Package CIDs and their metadata into a single struct
	return &CIDPayload{
		BlockHash:       payload.BlockHash,
		BlockNumber:     payload.BlockNumber.String(),
		HeaderCID:       headerCid,
		UncleCIDs:       uncleCids,
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
	// Keep receipts associated with their transaction
	mappedRctCids := make(map[common.Hash]*ReceiptMetaData, len(receiptsCids))
	for i, rct := range receipts {
		mappedRctCids[rct.TxHash] = receiptMeta[i]
		mappedRctCids[rct.TxHash].CID = receiptsCids[i]
	}
	return mappedRctCids, nil
}

func (pub *Publisher) publishStateNodes(stateNodes map[common.Hash]StateNode) (map[common.Hash]StateNodeCID, error) {
	stateNodeCids := make(map[common.Hash]StateNodeCID)
	for addrKey, node := range stateNodes {
		stateNodeCid, err := pub.StatePutter.DagPut(node.Value)
		if err != nil {
			return nil, err
		}
		if len(stateNodeCid) != 1 {
			return nil, errors.New("single CID expected to be returned for state leaf")
		}
		stateNodeCids[addrKey] = StateNodeCID{
			CID:  stateNodeCid[0],
			Leaf: node.Leaf,
		}
	}
	return stateNodeCids, nil
}

func (pub *Publisher) publishStorageNodes(storageNodes map[common.Hash][]StorageNode) (map[common.Hash][]StorageNodeCID, error) {
	storageLeafCids := make(map[common.Hash][]StorageNodeCID)
	for addrKey, storageTrie := range storageNodes {
		storageLeafCids[addrKey] = make([]StorageNodeCID, 0, len(storageTrie))
		for _, node := range storageTrie {
			storageNodeCid, err := pub.StoragePutter.DagPut(node.Value)
			if err != nil {
				return nil, err
			}
			if len(storageNodeCid) != 1 {
				return nil, errors.New("single CID expected to be returned for storage leaf")
			}
			storageLeafCids[addrKey] = append(storageLeafCids[addrKey], StorageNodeCID{
				Key:  node.Key.Hex(),
				CID:  storageNodeCid[0],
				Leaf: node.Leaf,
			})
		}
	}
	return storageLeafCids, nil
}
