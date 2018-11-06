// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"math/rand"
	"time"
)

type GenericModel struct{}
type GenericEntity struct{}

var startingBlockNumber = rand.Int63()
var topic = "0x" + randomString(64)
var address = "0x" + randomString(38)

var GenericTestLogs = []types.Log{{
	Address:     common.HexToAddress(address),
	Topics:      []common.Hash{common.HexToHash(topic)},
	BlockNumber: uint64(startingBlockNumber),
}}

var GenericTestConfig = shared.SingleTransformerConfig{
	TransformerName:     "generic-test-transformer",
	ContractAddresses:   []string{address},
	ContractAbi:         randomString(100),
	Topic:               topic,
	StartingBlockNumber: startingBlockNumber,
	EndingBlockNumber:   startingBlockNumber + 1,
}

func randomString(length int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyz1234567890"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
