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
	"math/big"

	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ConvertedPayload is a custom type which packages raw ETH data for publishing to IPFS and filtering to subscribers
// Returned by PayloadConverter
// Passed to IPLDPublisher and ResponseFilterer
type ConvertedPayload struct {
	TotalDifficulty *big.Int
	Block           *types.Block
	TxMetaData      []TxModel
	Receipts        types.Receipts
	ReceiptMetaData []ReceiptModel
	StateNodes      []TrieNode
	StorageNodes    map[common.Hash][]TrieNode
}

// Height satisfies the StreamedIPLDs interface
func (i ConvertedPayload) Height() int64 {
	return i.Block.Number().Int64()
}

// Trie struct used to flag node as leaf or not
type TrieNode struct {
	Path    []byte
	LeafKey common.Hash
	Value   []byte
	Type    statediff.NodeType
}

// CIDPayload is a struct to hold all the CIDs and their associated meta data for indexing in Postgres
// Returned by IPLDPublisher
// Passed to CIDIndexer
type CIDPayload struct {
	HeaderCID       HeaderModel
	UncleCIDs       []UncleModel
	TransactionCIDs []TxModel
	ReceiptCIDs     map[common.Hash]ReceiptModel
	StateNodeCIDs   []StateNodeModel
	StorageNodeCIDs map[common.Hash][]StorageNodeModel
}

// CIDWrapper is used to direct fetching of IPLDs from IPFS
// Returned by CIDRetriever
// Passed to IPLDFetcher
type CIDWrapper struct {
	BlockNumber  *big.Int
	Header       HeaderModel
	Uncles       []UncleModel
	Transactions []TxModel
	Receipts     []ReceiptModel
	StateNodes   []StateNodeModel
	StorageNodes []StorageNodeWithStateKeyModel
}

// IPLDs is used to package raw IPLD block data fetched from IPFS and returned by the server
// Returned by IPLDFetcher and ResponseFilterer
type IPLDs struct {
	BlockNumber     *big.Int
	TotalDifficulty *big.Int
	Header          ipfs.BlockModel
	Uncles          []ipfs.BlockModel
	Transactions    []ipfs.BlockModel
	Receipts        []ipfs.BlockModel
	StateNodes      []StateNode
	StorageNodes    []StorageNode
}

// Height satisfies the StreamedIPLDs interface
func (i IPLDs) Height() int64 {
	return i.BlockNumber.Int64()
}

type StateNode struct {
	Type         statediff.NodeType
	StateLeafKey common.Hash
	Path         []byte
	IPLD         ipfs.BlockModel
}

type StorageNode struct {
	Type           statediff.NodeType
	StateLeafKey   common.Hash
	StorageLeafKey common.Hash
	Path           []byte
	IPLD           ipfs.BlockModel
}
