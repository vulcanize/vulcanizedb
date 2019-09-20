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

package fakes

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var (
	// fakeBloom var should mean that the bloom tells us that the log might contain "FakeTopic"
	// A bunch of tests (like event_watcher_test) use a transformer config that sets "FakeTopic" as the topic0
	// We want to make sure those tests still pass (fetch logs needs to get called, so the bloom filter has to have bits set for this string)
	fakeBloom     = GetFakeBloom([]string{"FakeTopic"})
	FakeAddress   = common.HexToAddress("0x1234567890abcdef")
	FakeError     = errors.New("failed")
	FakeHash      = common.BytesToHash([]byte{1, 2, 3, 4, 5})
	fakeTimestamp = int64(111111111)
)

var rawFakeHeader, _ = json.Marshal(types.Header{})
var FakeHeader = core.Header{
	Bloom:     fakeBloom,
	Hash:      FakeHash.String(),
	Raw:       rawFakeHeader,
	Timestamp: strconv.FormatInt(fakeTimestamp, 10),
}

func GetFakeHeader(blockNumber int64) core.Header {
	return GetFakeHeaderWithTimestamp(fakeTimestamp, blockNumber)
}

func GetFakeHeaderWithTimestamp(timestamp, blockNumber int64) core.Header {
	return core.Header{
		Bloom:       fakeBloom,
		Hash:        FakeHash.String(),
		BlockNumber: blockNumber,
		Raw:         rawFakeHeader,
		Timestamp:   strconv.FormatInt(timestamp, 10),
	}
}

var fakeTransaction types.Transaction
var rawTransaction bytes.Buffer
var _ = fakeTransaction.EncodeRLP(&rawTransaction)
var FakeTransaction = core.TransactionModel{
	Data:     []byte{},
	From:     "",
	GasLimit: 0,
	GasPrice: 0,
	Hash:     "",
	Nonce:    0,
	Raw:      rawTransaction.Bytes(),
	Receipt:  core.Receipt{},
	To:       "",
	TxIndex:  0,
	Value:    "0",
}

func GetFakeTransaction(hash string, receipt core.Receipt) core.TransactionModel {
	gethTransaction := types.Transaction{}
	var raw bytes.Buffer
	err := gethTransaction.EncodeRLP(&raw)
	if err != nil {
		panic("failed to marshal transaction while creating test fake")
	}
	return core.TransactionModel{
		Data:     []byte{},
		From:     "",
		GasLimit: 0,
		GasPrice: 0,
		Hash:     hash,
		Nonce:    0,
		Raw:      raw.Bytes(),
		Receipt:  receipt,
		To:       "",
		TxIndex:  0,
		Value:    "0",
	}
}

func GetFakeUncle(hash, reward string) core.Uncle {
	return core.Uncle{
		Miner:     FakeAddress.String(),
		Hash:      hash,
		Reward:    reward,
		Raw:       rawFakeHeader,
		Timestamp: strconv.FormatInt(fakeTimestamp, 10),
	}
}

// this takes a string array and convert it to the correct type for a topic0
// then adds it to the blooom filter, so we can make sure that bloom.Test will return a positive result for the given strings
func GetFakeBloom(positive []string) []byte {
	var bloom types.Bloom

	for _, data := range positive {
		topicHash := common.HexToHash(data)
		bloom.Add(new(big.Int).SetBytes(topicHash.Bytes()))
	}
	return bloom.Bytes()
}

func GetFakeHeaderWithPositiveBloom(positive []string) core.Header {
	return core.Header{
		Bloom: GetFakeBloom(positive),
		Hash:      FakeHeader.Hash,
		Raw:       FakeHeader.Raw,
		Timestamp: FakeHeader.Timestamp,
	}
}
