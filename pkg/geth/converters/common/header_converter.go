package common

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type HeaderConverter struct{}

func (converter HeaderConverter) Convert(gethHeader *types.Header, blockHash string) (core.Header, error) {
	rawHeader, err := json.Marshal(gethHeader)
	if err != nil {
		panic(err)
	}
	coreHeader := core.Header{
		Hash:        blockHash,
		BlockNumber: gethHeader.Number.Int64(),
		Raw:         rawHeader,
		Timestamp:   gethHeader.Time.String(),
	}
	return coreHeader, nil
}
