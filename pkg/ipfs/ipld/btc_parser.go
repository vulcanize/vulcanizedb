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
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	node "github.com/ipfs/go-ipld-format"
)

// FromHeaderAndTxs takes a block header and txs and processes it
// to return it a set of IPLD nodes for further processing.
func FromHeaderAndTxs(header *wire.BlockHeader, txs []*btcutil.Tx) (*BtcHeader, []*BtcTx, []*BtcTxTrie, error) {
	var txNodes []*BtcTx
	for _, tx := range txs {
		txNode, err := NewBtcTx(tx.MsgTx())
		if err != nil {
			return nil, nil, nil, err
		}
		txNodes = append(txNodes, txNode)
	}
	txTrie, err := mkMerkleTree(txNodes)
	if err != nil {
		return nil, nil, nil, err
	}
	headerNode, err := NewBtcHeader(header)
	return headerNode, txNodes, txTrie, err
}

func mkMerkleTree(txs []*BtcTx) ([]*BtcTxTrie, error) {
	layer := make([]node.Node, len(txs))
	for i, tx := range txs {
		layer[i] = tx
	}
	var out []*BtcTxTrie
	var next []node.Node
	for len(layer) > 1 {
		if len(layer)%2 != 0 {
			layer = append(layer, layer[len(layer)-1])
		}
		for i := 0; i < len(layer)/2; i++ {
			var left, right node.Node
			left = layer[i*2]
			right = layer[(i*2)+1]

			t := &BtcTxTrie{
				Left:  &node.Link{Cid: left.Cid()},
				Right: &node.Link{Cid: right.Cid()},
			}

			out = append(out, t)
			next = append(next, t)
		}

		layer = next
		next = nil
	}

	return out, nil
}
