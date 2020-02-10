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

package ipld

import (
	"fmt"

	"github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

// EthStateTrie (eth-state-trie, codec 0x96), represents
// a node from the state trie in ethereum.
type EthStateTrie struct {
	cid     cid.Cid
	rawdata []byte
}

// Static (compile time) check that EthStateTrie satisfies the node.Node interface.
var _ node.Node = (*EthStateTrie)(nil)

/*
  INPUT
*/

// FromStateTrieRLP takes the RLP bytes of an ethereum
// state trie node to return it as an IPLD node for further processing.
func FromStateTrieRLP(stateNodeRLP []byte) (*EthStateTrie, error) {
	c, err := rawdataToCid(MEthStateTrie, stateNodeRLP, mh.KECCAK_256)
	if err != nil {
		return nil, err
	}
	return &EthStateTrie{
		cid:     c,
		rawdata: stateNodeRLP,
	}, nil
}

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the state trie node.
func (st *EthStateTrie) RawData() []byte {
	return st.rawdata
}

// Cid returns the cid of the state trie node.
func (st *EthStateTrie) Cid() cid.Cid {
	return st.cid
}

// String is a helper for output
func (st *EthStateTrie) String() string {
	return fmt.Sprintf("<EthereumStateTrie %s>", st.cid)
}

// Copy will go away. It is here to comply with the Node interface.
func (*EthStateTrie) Copy() node.Node {
	panic("implement me")
}

func (*EthStateTrie) Links() []*node.Link {
	panic("implement me")
}

func (*EthStateTrie) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (*EthStateTrie) ResolveLink(path []string) (*node.Link, []string, error) {
	panic("implement me")
}

func (*EthStateTrie) Size() (uint64, error) {
	panic("implement me")
}

func (*EthStateTrie) Stat() (*node.NodeStat, error) {
	panic("implement me")
}

func (*EthStateTrie) Tree(path string, depth int) []string {
	panic("implement me")
}

// Loggable returns in a map the type of IPLD Link.
func (st *EthStateTrie) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth-state-trie",
	}
}
