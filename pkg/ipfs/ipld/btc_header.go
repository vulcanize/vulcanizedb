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
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

type BtcHeader struct {
	*wire.BlockHeader

	rawdata []byte
	cid     cid.Cid
}

// Static (compile time) check that BtcBtcHeader satisfies the node.Node interface.
var _ node.Node = (*BtcHeader)(nil)

/*
  INPUT
*/

// NewBtcHeader converts a *wire.Header into an BtcHeader IPLD node
func NewBtcHeader(header *wire.BlockHeader) (*BtcHeader, error) {
	w := bytes.NewBuffer(make([]byte, 0, 80))
	if err := header.Serialize(w); err != nil {
		return nil, err
	}
	rawdata := w.Bytes()
	c, err := RawdataToCid(MBitcoinHeader, rawdata, mh.DBL_SHA2_256)
	if err != nil {
		return nil, err
	}
	return &BtcHeader{
		BlockHeader: header,
		cid:         c,
		rawdata:     rawdata,
	}, nil
}

/*
   Block INTERFACE
*/

func (b *BtcHeader) Cid() cid.Cid {
	return b.cid
}

func (b *BtcHeader) RawData() []byte {
	return b.rawdata
}

func (b *BtcHeader) String() string {
	return fmt.Sprintf("<BtcHeader %s>", b.cid)
}

func (b *BtcHeader) Loggable() map[string]interface{} {
	// TODO: more helpful info here
	return map[string]interface{}{
		"type": "bitcoin_block",
	}
}

/*
   Node INTERFACE
*/

func (b *BtcHeader) Links() []*node.Link {
	return []*node.Link{
		{
			Name: "tx",
			Cid:  sha256ToCid(MBitcoinTx, b.MerkleRoot.CloneBytes()),
		},
		{
			Name: "parent",
			Cid:  sha256ToCid(MBitcoinHeader, b.PrevBlock.CloneBytes()),
		},
	}
}

// Resolve attempts to traverse a path through this block.
func (b *BtcHeader) Resolve(path []string) (interface{}, []string, error) {
	if len(path) == 0 {
		return nil, nil, fmt.Errorf("zero length path")
	}
	switch path[0] {
	case "version":
		return b.Version, path[1:], nil
	case "timestamp":
		return b.Timestamp, path[1:], nil
	case "bits":
		return b.Bits, path[1:], nil
	case "nonce":
		return b.Nonce, path[1:], nil
	case "parent":
		return &node.Link{Cid: sha256ToCid(MBitcoinHeader, b.PrevBlock.CloneBytes())}, path[1:], nil
	case "tx":
		return &node.Link{Cid: sha256ToCid(MBitcoinTx, b.MerkleRoot.CloneBytes())}, path[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

// ResolveLink is a helper function that allows easier traversal of links through blocks
func (b *BtcHeader) ResolveLink(path []string) (*node.Link, []string, error) {
	out, rest, err := b.Resolve(path)
	if err != nil {
		return nil, nil, err
	}

	lnk, ok := out.(*node.Link)
	if !ok {
		return nil, nil, fmt.Errorf("object at path was not a link")
	}

	return lnk, rest, nil
}

func cidToHash(c cid.Cid) []byte {
	h := []byte(c.Hash())
	return h[len(h)-32:]
}

func hashToCid(hv []byte, t uint64) cid.Cid {
	h, _ := mh.Encode(hv, mh.DBL_SHA2_256)
	return cid.NewCidV1(t, h)
}

func (b *BtcHeader) Size() (uint64, error) {
	return uint64(len(b.rawdata)), nil
}

func (b *BtcHeader) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

func (b *BtcHeader) Tree(p string, depth int) []string {
	// TODO: this isnt a correct implementation yet
	return []string{"difficulty", "nonce", "version", "timestamp", "tx", "parent"}
}

func (b *BtcHeader) BTCSha() []byte {
	blkmh, _ := mh.Sum(b.rawdata, mh.DBL_SHA2_256, -1)
	return blkmh[2:]
}

func (b *BtcHeader) HexHash() string {
	return hex.EncodeToString(revString(b.BTCSha()))
}

func (b *BtcHeader) Copy() node.Node {
	nb := *b // cheating shallow copy
	return &nb
}

func revString(s []byte) []byte {
	b := make([]byte, len(s))
	for i, v := range []byte(s) {
		b[len(b)-(i+1)] = v
	}
	return b
}
