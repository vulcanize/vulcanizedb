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

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_state_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_storage_trie"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

type Publisher interface {
	Publish(payload *IPLDPayload) (*CIDPayload, error)
}

type IPLDPublisher struct {
	Node *ipfs.IPFS
	HeaderPutter *eth_block_header.BlockHeaderDagPutter
	TransactionPutter *eth_block_transactions.BlockTransactionsDagPutter
	ReceiptPutter *eth_block_receipts.EthBlockReceiptDagPutter
	StatePutter *eth_state_trie.StateTrieDagPutter
	StoragePutter *eth_storage_trie.StorageTrieDagPutter
}

type CIDPayload struct {
	BlockNumber string
	BlockHash  string
	HeaderCID  string
	TransactionCIDs map[common.Hash]string
	ReceiptCIDs map[common.Hash]string
	StateLeafCIDs  map[common.Hash]string
	StorageLeafCIDs map[common.Hash]map[common.Hash]string
}

func NewIPLDPublisher(ipfsPath string) (*IPLDPublisher, error) {
	node, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		return nil, err
	}
	decoder := rlp.RlpDecoder{}
	return &IPLDPublisher{
		Node: node,
		HeaderPutter: eth_block_header.NewBlockHeaderDagPutter(node, decoder),
		TransactionPutter: eth_block_transactions.NewBlockTransactionsDagPutter(node),
		ReceiptPutter: eth_block_receipts.NewEthBlockReceiptDagPutter(node),
		StatePutter: eth_state_trie.NewStateTrieDagPutter(node),
		StoragePutter: eth_storage_trie.NewStorageTrieDagPutter(node),
	}, nil
}

func (pub *IPLDPublisher) Publish(payload *IPLDPayload) (*CIDPayload, error) {
	// Process and publish headers
	headerCids, err := pub.HeaderPutter.DagPut(payload.HeaderRLP)
	if err != nil {
		return nil, err
	}
	if len(headerCids) != 1 {
		return nil, errors.New("single CID expected to be returned for header")
	}

	// Process and publish transactions
	transactionCids, err := pub.TransactionPutter.DagPut(payload.BlockBody)
	if err != nil {
		return nil, err
	}
	if len(transactionCids) != len(payload.BlockBody.Transactions) {
		return nil, errors.New("expected one CID for each transaction")
	}
	trxCids := make(map[common.Hash]string, len(transactionCids))
	for i, trx := range payload.BlockBody.Transactions {
		trxCids[trx.Hash()] = transactionCids[i]
	}

	// Process and publish receipts
	receiptsCids, err := pub.ReceiptPutter.DagPut(payload.Receipts)
	if err != nil {
		return nil, err
	}
	if len(receiptsCids) != len(payload.Receipts) {
		return nil, errors.New("expected one CID for each receipt")
	}
	rctCids := make(map[common.Hash]string, len(receiptsCids))
	for i, rct := range payload.Receipts {
		rctCids[rct.TxHash] = receiptsCids[i]
	}

	// Process and publish state leafs
	stateLeafCids := make(map[common.Hash]string)
	for addr, leaf := range payload.StateLeafs {
		stateLeafCid, err := pub.StatePutter.DagPut(leaf)
		if err != nil {
			return nil, err
		}
		if len(stateLeafCid) != 1 {
			return nil, errors.New("single CID expected to be returned for state leaf")
		}
		stateLeafCids[addr] = stateLeafCid[0]
	}

	// Process and publish storage leafs
	storageLeafCids := make(map[common.Hash]map[common.Hash]string)
	for addr, storageTrie := range payload.StorageLeafs {
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

	return &CIDPayload{
		BlockHash: payload.BlockHash.Hex(),
		BlockNumber: payload.BlockNumber.String(),
		HeaderCID: headerCids[0],
		TransactionCIDs: trxCids,
		ReceiptCIDs: rctCids,
		StateLeafCIDs: stateLeafCids,
		StorageLeafCIDs: storageLeafCids,
	}, nil
}