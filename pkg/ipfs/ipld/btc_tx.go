package ipld

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/wire"
	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

type BtcTx struct {
	*wire.MsgTx

	rawdata []byte
	cid     cid.Cid
}

// Static (compile time) check that BtcBtcHeader satisfies the node.Node interface.
var _ node.Node = (*BtcTx)(nil)

/*
  INPUT
*/

// NewBtcTx converts a *wire.MsgTx into an BtcTx IPLD node
func NewBtcTx(tx *wire.MsgTx) (*BtcTx, error) {
	w := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(w); err != nil {
		return nil, err
	}
	rawdata := w.Bytes()
	c, err := RawdataToCid(MBitcoinTx, rawdata, mh.DBL_SHA2_256)
	if err != nil {
		return nil, err
	}
	return &BtcTx{
		MsgTx:   tx,
		cid:     c,
		rawdata: rawdata,
	}, nil
}

/*
   Block INTERFACE
*/

func (t *BtcTx) Cid() cid.Cid {
	return t.cid
}

func (t *BtcTx) RawData() []byte {
	return t.rawdata
}

func (t *BtcTx) String() string {
	return fmt.Sprintf("<BtcTx %s>", t.cid)
}

func (t *BtcTx) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "bitcoinTx",
	}
}

/*
   Node INTERFACE
*/

func (t *BtcTx) Links() []*node.Link {
	var out []*node.Link
	for i, in := range t.MsgTx.TxIn {
		lnk := &node.Link{Cid: sha256ToCid(MBitcoinTx, in.PreviousOutPoint.Hash.CloneBytes())}
		lnk.Name = fmt.Sprintf("inputs/%d/prevTx", i)
		out = append(out, lnk)
	}
	return out
}

func (t *BtcTx) Resolve(path []string) (interface{}, []string, error) {
	switch path[0] {
	case "version":
		return t.Version, path[1:], nil
	case "lockTime":
		return t.LockTime, path[1:], nil
	case "inputs":
		if len(path) == 1 {
			return t.MsgTx.TxIn, nil, nil
		}

		index, err := strconv.Atoi(path[1])
		if err != nil {
			return nil, nil, err
		}

		if index >= len(t.MsgTx.TxIn) || index < 0 {
			return nil, nil, fmt.Errorf("index out of range")
		}

		inp := t.MsgTx.TxIn[index]
		if len(path) == 2 {
			return inp, nil, nil
		}

		switch path[2] {
		case "prevTx":
			return &node.Link{Cid: sha256ToCid(MBitcoinTx, inp.PreviousOutPoint.Hash.CloneBytes())}, path[3:], nil
		case "seqNo":
			return inp.Sequence, path[3:], nil
		case "script":
			return inp.SignatureScript, path[3:], nil
		default:
			return nil, nil, fmt.Errorf("no such link")
		}
	case "outputs":
		if len(path) == 1 {
			return t.TxOut, nil, nil
		}

		index, err := strconv.Atoi(path[1])
		if err != nil {
			return nil, nil, err
		}

		if index >= len(t.TxOut) || index < 0 {
			return nil, nil, fmt.Errorf("index out of range")
		}

		outp := t.TxOut[index]
		if len(path) == 2 {
			return outp, path[2:], nil
		}

		switch path[2] {
		case "value":
			return outp.Value, path[3:], nil
		case "script":
			/*
				if outp.Script[0] == 0x6a { // OP_RETURN
					c, err := cid.Decode(string(outp.Script[1:]))
					if err == nil {
						return &node.Link{Cid: c}, path[3:], nil
					}
				}
			*/
			return outp.PkScript, path[3:], nil
		default:
			return nil, nil, fmt.Errorf("no such link")
		}
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

func (t *BtcTx) ResolveLink(path []string) (*node.Link, []string, error) {
	i, rest, err := t.Resolve(path)
	if err != nil {
		return nil, rest, err
	}

	lnk, ok := i.(*node.Link)
	if !ok {
		return nil, nil, fmt.Errorf("value was not a link")
	}

	return lnk, rest, nil
}

func (t *BtcTx) Size() (uint64, error) {
	return uint64(len(t.RawData())), nil
}

func (t *BtcTx) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

func (t *BtcTx) Copy() node.Node {
	nt := *t // cheating shallow copy
	return &nt
}

func (t *BtcTx) Tree(p string, depth int) []string {
	if depth == 0 {
		return nil
	}

	switch p {
	case "inputs":
		return t.treeInputs(nil, depth+1)
	case "outputs":
		return t.treeOutputs(nil, depth+1)
	case "":
		out := []string{"version", "timeLock", "inputs", "outputs"}
		out = t.treeInputs(out, depth)
		out = t.treeOutputs(out, depth)
		return out
	default:
		return nil
	}
}

func (t *BtcTx) treeInputs(out []string, depth int) []string {
	if depth < 2 {
		return out
	}

	for i := range t.TxIn {
		inp := "inputs/" + fmt.Sprint(i)
		out = append(out, inp)
		if depth > 2 {
			out = append(out, inp+"/prevTx", inp+"/seqNo", inp+"/script")
		}
	}
	return out
}

func (t *BtcTx) treeOutputs(out []string, depth int) []string {
	if depth < 2 {
		return out
	}

	for i := range t.TxOut {
		o := "outputs/" + fmt.Sprint(i)
		out = append(out, o)
		if depth > 2 {
			out = append(out, o+"/script", o+"/value")
		}
	}
	return out
}

func (t *BtcTx) BTCSha() []byte {
	mh, _ := mh.Sum(t.RawData(), mh.DBL_SHA2_256, -1)
	return []byte(mh[2:])
}

func (t *BtcTx) HexHash() string {
	return hex.EncodeToString(revString(t.BTCSha()))
}
