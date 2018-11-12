package fakes

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var (
	FakeError = errors.New("failed")
	FakeHash  = common.BytesToHash([]byte{1, 2, 3, 4, 5})
)

var rawFakeHeader, _ = json.Marshal(types.Header{})
var FakeHeader = core.Header{
	Hash:      FakeHash.String(),
	Raw:       rawFakeHeader,
	Timestamp: "111111111",
}

func GetFakeHeader(blockNumber int64) core.Header {
	return core.Header{
		Hash:        FakeHash.String(),
		BlockNumber: blockNumber,
		Raw:         rawFakeHeader,
		Timestamp:   "111111111",
	}
}
