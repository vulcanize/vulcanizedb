package common

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type HeaderConverter struct{}

func (converter HeaderConverter) Convert(gethHeader *types.Header) (core.Header, error) {
	writer := new(bytes.Buffer)
	err := rlp.Encode(writer, &gethHeader)
	if err != nil {
		panic(err)
	}
	coreHeader := core.Header{
		Hash:        gethHeader.Hash().Hex(),
		BlockNumber: gethHeader.Number.Int64(),
		Raw:         writer.Bytes(),
	}
	return coreHeader, nil
}
