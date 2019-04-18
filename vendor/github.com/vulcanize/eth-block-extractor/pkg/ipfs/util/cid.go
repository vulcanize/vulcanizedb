package util

import (
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

func RawToCid(codec uint64, raw []byte) (cid.Cid, error) {
	c, err := cid.Prefix{
		Codec:    codec,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(raw)
	if err != nil {
		return cid.Cid{}, err
	}
	return c, nil
}
