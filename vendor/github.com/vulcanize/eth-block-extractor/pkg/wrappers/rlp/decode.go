package rlp

import (
	"bytes"
	"github.com/ethereum/go-ethereum/rlp"
)

type Decoder interface {
	Decode(raw []byte, out interface{}) error
}

type RlpDecoder struct{}

func (RlpDecoder) Decode(raw []byte, out interface{}) error {
	return rlp.Decode(bytes.NewBuffer(raw), out)
}
