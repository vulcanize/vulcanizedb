package fakes

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"strconv"
)

var (
	FakeError     = errors.New("failed")
	FakeHash      = common.BytesToHash([]byte{1, 2, 3, 4, 5})
	fakeTimestamp = int64(111111111)
)

var rawFakeHeader, _ = json.Marshal(types.Header{})
var FakeHeader = core.Header{
	Hash:      FakeHash.String(),
	Raw:       rawFakeHeader,
	Timestamp: strconv.FormatInt(fakeTimestamp, 10),
}

func GetFakeHeader(blockNumber int64) core.Header {
	return core.Header{
		Hash:        FakeHash.String(),
		BlockNumber: blockNumber,
		Raw:         rawFakeHeader,
		Timestamp:   strconv.FormatInt(fakeTimestamp, 10),
	}
}

var FakeHeaderTic = fakeTimestamp + constants.TTL
