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
	stateLeafCids, err := pub.publishStateLeafs(payload.StateLeafs)
	if err != nil {
		return nil, err
	}

	// Process and publish storage leafs
	storageLeafCids, err := pub.publishStorageLeafs(payload.StorageLeafs)
	if err != nil {
		return nil, err
	}

	// Package CIDs into a single struct
	return &CIDPayload{
		BlockHash:       payload.BlockHash.Hex(),
		BlockNumber:     payload.BlockNumber.String(),
		HeaderCID:       headerCid,
		UncleCIDS:       uncleCids,
		TransactionCIDs: transactionCids,
		ReceiptCIDs:     receiptsCids,
		StateLeafCIDs:   stateLeafCids,
		StorageLeafCIDs: storageLeafCids,
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
	mappedRctCids := make(map[common.Hash]*ReceiptMetaData, len(receiptsCids))
	for i, rct := range receipts {
		mappedRctCids[rct.TxHash] = receiptMeta[i]
		mappedRctCids[rct.TxHash].CID = receiptsCids[i]
	}
	return mappedRctCids, nil
}

func (pub *Publisher) publishStateLeafs(stateLeafs map[common.Hash][]byte) (map[common.Hash]string, error) {
	stateLeafCids := make(map[common.Hash]string)
	for addr, leaf := range stateLeafs {
		stateLeafCid, err := pub.StatePutter.DagPut(leaf)
		if err != nil {
			return nil, err
		}
		if len(stateLeafCid) != 1 {
			return nil, errors.New("single CID expected to be returned for state leaf")
		}
		stateLeafCids[addr] = stateLeafCid[0]
	}
	return stateLeafCids, nil
}

func (pub *Publisher) publishStorageLeafs(storageLeafs map[common.Hash]map[common.Hash][]byte) (map[common.Hash]map[common.Hash]string, error) {
	storageLeafCids := make(map[common.Hash]map[common.Hash]string)
	for addr, storageTrie := range storageLeafs {
		storageLeafCids[addr] = make(map[common.Hash]string)
		for key, leaf := range storageTrie {
			storageLeafCid, err := pub.StoragePutter.DagPut(leaf)
			if err != nil {
				return nil, err
			}
			if len(storageLeafCid) != 1 {
				return nil, errors.New("single CID expected to be returned for storage leaf")
			}
			storageLeafCids[addr][key] = storageLeafCid[0]
		}
	}
	return storageLeafCids, nil
}
