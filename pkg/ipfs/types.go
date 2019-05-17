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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//
type Subscription struct {
	PayloadChan chan<- ResponsePayload
	QuitChan    chan<- bool
}

type ResponsePayload struct {
	HeadersRlp      [][]byte `json:"headersRlp"`
	UnclesRlp       [][]byte `json:"unclesRlp"`
	TransactionsRlp [][]byte `json:"transactionsRlp"`
	ReceiptsRlp     [][]byte `json:"receiptsRlp"`
	StateNodesRlp   [][]byte `json:"stateNodesRlp"`
	StorageNodesRlp [][]byte `json:"storageNodesRlp"`

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

type StateNode struct {
	Value []byte
	Leaf  bool
}

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

type StateNodeCID struct {
	CID  string
	Leaf bool
}

type StorageNodeCID struct {
	Key  common.Hash
	CID  string
	Leaf bool
}

// ReceiptMetaData wraps some additional data around our receipt CIDs for indexing
type ReceiptMetaData struct {
	CID     string
	Topic0s []string
}

// TrxMetaData wraps some additional data around our transaction CID for indexing
type TrxMetaData struct {
	CID  string
	To   string
	From string
}

// Params are set by the client to tell the server how to filter that is fed into their subscription
type Params struct {
	HeaderFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64 // set to 0 or a negative value to have no ending block
		Uncles        bool
	}
	TrxFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Src           string
		Dst           string
	}
	ReceiptFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Topic0s       []string
	}
	StateFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Address       string // is converted to state key by taking its keccak256 hash
		LeafsOnly     bool
	}
	StorageFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Address       string
		StorageKey    string
		LeafsOnly     bool
	}
}
