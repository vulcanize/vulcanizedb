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
	"encoding/json"
	"math/big"

	"github.com/ipfs/go-block-format"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Subscription holds the information for an individual client subscription
type Subscription struct {
	PayloadChan   chan<- ResponsePayload
	QuitChan      chan<- bool
	StreamFilters *StreamFilters
}

// ResponsePayload holds the data returned from the seed node to the requesting client
type ResponsePayload struct {
	HeadersRlp      [][]byte                               `json:"headersRlp"`
	UnclesRlp       [][]byte                               `json:"unclesRlp"`
	TransactionsRlp [][]byte                               `json:"transactionsRlp"`
	ReceiptsRlp     [][]byte                               `json:"receiptsRlp"`
	StateNodesRlp   map[common.Hash][]byte                 `json:"stateNodesRlp"`
	StorageNodesRlp map[common.Hash]map[common.Hash][]byte `json:"storageNodesRlp"`
	Err             error                                  `json:"error"`

	encoded []byte
	err     error
}

func (sd *ResponsePayload) ensureEncoded() {
	if sd.encoded == nil && sd.err == nil {
		sd.encoded, sd.err = json.Marshal(sd)
	}
}

// Length to implement Encoder interface for StateDiff
func (sd *ResponsePayload) Length() int {
	sd.ensureEncoded()
	return len(sd.encoded)
}

// Encode to implement Encoder interface for StateDiff
func (sd *ResponsePayload) Encode() ([]byte, error) {
	sd.ensureEncoded()
	return sd.encoded, sd.err
}

// CidWrapper is used to package CIDs retrieved from the local Postgres cache
type CidWrapper struct {
	BlockNumber  int64
	Headers      []string
	Transactions []string
	Receipts     []string
	StateNodes   []StateNodeCID
	StorageNodes []StorageNodeCID
}

// IpldWrapper is used to package raw IPLD block data for resolution
type IpldWrapper struct {
	Headers      []blocks.Block
	Transactions []blocks.Block
	Receipts     []blocks.Block
	StateNodes   map[common.Hash]blocks.Block
	StorageNodes map[common.Hash]map[common.Hash]blocks.Block
}

// IPLDPayload is a custom type which packages ETH data for the IPFS publisher
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
	UncleCIDS       map[common.Hash]string
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
	CID     string
	Topic0s []string
}

// TrxMetaData wraps some additional data around our transaction CID for indexing
type TrxMetaData struct {
	CID string
	Src string
	Dst string
}

// StreamFilters are defined by the client to specifiy which data to receive from the seed node
type StreamFilters struct {
	BackFill      bool
	BackFillOnly  bool
	StartingBlock int64
	EndingBlock   int64 // set to 0 or a negative value to have no ending block
	HeaderFilter  struct {
		Off       bool
		FinalOnly bool
	}
	TrxFilter struct {
		Off bool
		Src []string
		Dst []string
	}
	ReceiptFilter struct {
		Off     bool
		Topic0s []string
	}
	StateFilter struct {
		Off               bool
		Addresses         []string // is converted to state key by taking its keccak256 hash
		IntermediateNodes bool
	}
	StorageFilter struct {
		Off               bool
		Addresses         []string
		StorageKeys       []string
		IntermediateNodes bool
	}
}
