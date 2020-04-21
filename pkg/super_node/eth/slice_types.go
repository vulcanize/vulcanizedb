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
)

// GetSliceInput holds the arguments to the eth_getSlice method
type GetSliceInput struct {
	Path    string
	Depth   int
	Root    common.Hash
	Storage bool
}

// GetSliceResponse holds response for the eth_getSlice method
type GetSliceResponse struct {
	SliceID   string                             `json:"slice-id"`
	MetaData  GetSliceResponseMetadata           `json:"metadata"`
	TrieNodes GetSliceResponseTrieNodes          `json:"trie-nodes"`
	Leaves    map[string]GetSliceResponseAccount `json:"leaves"` // we won't be using addresses, but keccak256(address)
}

func (sr *GetSliceResponse) init(path string, depth int, root common.Hash) {
	sr.SliceID = fmt.Sprintf("%s-%d-%s", path, depth, root.String())
	sr.MetaData = GetSliceResponseMetadata{
		NodeStats: make(map[string]string, 0),
		TimeStats: make(map[string]string, 0),
	}
	sr.Leaves = make(map[string]GetSliceResponseAccount)
	sr.TrieNodes = GetSliceResponseTrieNodes{
		Stem:  make(map[string]string),
		Head:  make(map[string]string),
		Slice: make(map[string]string),
	}
}

type GetSliceResponseMetadata struct {
	TimeStats map[string]string `json:"time-ms"`    // stem, state, storage (one by one)
	NodeStats map[string]string `json:"trie-nodes"` // total, leaves, smart contracts
}

type GetSliceResponseTrieNodes struct {
	Stem  map[string]string `json:"stem"`
	Head  map[string]string `json:"head"`
	Slice map[string]string `json:"slice-nodes"`
}

type GetSliceResponseAccount struct {
	StorageRoot string `json:"storage-root"`
	EVMCode     string `json:"evm-code"`
}

type nodeDBResponse struct {
	CID      string `db:"cid"`
	NodeType int    `db:"node_type"`
}
