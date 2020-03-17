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

type BtcTxTrie struct {
	Left  *node.Link
	Right *node.Link
}

func (t *BtcTxTrie) BTCSha() []byte {
	return cidToHash(t.Cid())
}

func (t *BtcTxTrie) Cid() cid.Cid {
	h, _ := mh.Sum(t.RawData(), mh.DBL_SHA2_256, -1)
	return cid.NewCidV1(cid.BitcoinTx, h)
}

func (t *BtcTxTrie) Links() []*node.Link {
	return []*node.Link{t.Left, t.Right}
}

func (t *BtcTxTrie) RawData() []byte {
	out := make([]byte, 64)
	lbytes := cidToHash(t.Left.Cid)
	copy(out[:32], lbytes)

	rbytes := cidToHash(t.Right.Cid)
	copy(out[32:], rbytes)

	return out
}

func (t *BtcTxTrie) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "bitcoin_tx_tree",
	}
}

func (t *BtcTxTrie) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return nil, nil, fmt.Errorf("zero length path")
	}

	switch path[0] {
	case "0":
		return t.Left, path[1:], nil
	case "1":
		return t.Right, path[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

func (t *BtcTxTrie) Copy() node.Node {
	nt := *t
	return &nt
}

func (t *BtcTxTrie) ResolveLink(path []string) (*node.Link, []string, error) {
	out, rest, err := t.Resolve(path)
	if err != nil {
		return nil, nil, err
	}

	lnk, ok := out.(*node.Link)
	if ok {
		return lnk, rest, nil
	}

	return nil, nil, fmt.Errorf("path did not lead to link")
}

func (t *BtcTxTrie) Size() (uint64, error) {
	return uint64(len(t.RawData())), nil
}

func (t *BtcTxTrie) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

func (t *BtcTxTrie) String() string {
	return fmt.Sprintf("[bitcoin transaction tree]")
}

func (t *BtcTxTrie) Tree(p string, depth int) []string {
	return []string{"0", "1"}
}
