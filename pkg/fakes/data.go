package fakes

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var FakeError = errors.New("failed")

var rawFakeHeader, _ = json.Marshal(types.Header{})
var FakeHeader = core.Header{
	Raw:       rawFakeHeader,
	Timestamp: "111111111",
}

func GetFakeHeader(blockNumber int64) core.Header {
	return core.Header{
		BlockNumber: blockNumber,
		Raw:         rawFakeHeader,
		Timestamp:   "111111111",
	}
}
