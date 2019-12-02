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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ipfs/go-block-format"
)

// CIDWrapper is used to package CIDs retrieved from the local Postgres cache and direct fetching of IPLDs
type CIDWrapper struct {
	BlockNumber  *big.Int
	Headers      []string
	Uncles       []string
	Transactions []string
	Receipts     []string
	StateNodes   []StateNodeCID
	StorageNodes []StorageNodeCID
}

// IPLDWrapper is used to package raw IPLD block data fetched from IPFS
type IPLDWrapper struct {
	BlockNumber  *big.Int
	Headers      []blocks.Block
	Uncles       []blocks.Block
	Transactions []blocks.Block
	Receipts     []blocks.Block
	StateNodes   map[common.Hash]blocks.Block
	StorageNodes map[common.Hash]map[common.Hash]blocks.Block
}

// IPLDPayload is a custom type which packages raw ETH data for the IPFS publisher
type IPLDPayload struct {
	HeaderRLP       []byte
	BlockNumber     *big.Int
	BlockHash       common.Hash
	BlockBody       *types.Body
	TrxMetaData     []*TrxMetaData
	Receipts        types.Receipts
	ReceiptMetaData []*ReceiptMetaData
	StateNodes      map[common.Hash]StateNode
	StorageNodes    map[common.Hash][]StorageNode
}

// StateNode struct used to flag node as leaf or not
type StateNode struct {
	Value []byte
	Leaf  bool
}

// StorageNode struct used to flag node as leaf or not
type StorageNode struct {
	Key   common.Hash
	Value []byte
	Leaf  bool
}

// CIDPayload is a struct to hold all the CIDs and their meta data
type CIDPayload struct {
	BlockNumber     string
	BlockHash       common.Hash
	HeaderCID       string
	UncleCIDs       map[common.Hash]string
	TransactionCIDs map[common.Hash]*TrxMetaData
	ReceiptCIDs     map[common.Hash]*ReceiptMetaData
	StateNodeCIDs   map[common.Hash]StateNodeCID
	StorageNodeCIDs map[common.Hash][]StorageNodeCID
}

// StateNodeCID is used to associate a leaf flag with a state node cid
type StateNodeCID struct {
	CID  string
	Leaf bool
	Key  string `db:"state_key"`
}

// StorageNodeCID is used to associate a leaf flag with a storage node cid
type StorageNodeCID struct {
	Key      string `db:"storage_key"`
	CID      string
	Leaf     bool
	StateKey string `db:"state_key"`
}

// ReceiptMetaData wraps some additional data around our receipt CIDs for indexing
type ReceiptMetaData struct {
	CID             string
	Topic0s         []string
	ContractAddress string
}

// TrxMetaData wraps some additional data around our transaction CID for indexing
type TrxMetaData struct {
	CID string
	Src string
	Dst string
}
