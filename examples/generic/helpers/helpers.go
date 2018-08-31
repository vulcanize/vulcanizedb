// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func ConvertToLog(watchedEvent core.WatchedEvent) types.Log {
	allTopics := []string{watchedEvent.Topic0, watchedEvent.Topic1, watchedEvent.Topic2, watchedEvent.Topic3}
	var nonNilTopics []string
	for _, topic := range allTopics {
		if topic != "" {
			nonNilTopics = append(nonNilTopics, topic)
		}
	}
	return types.Log{
		Address:     common.HexToAddress(watchedEvent.Address),
		Topics:      createTopics(nonNilTopics...),
		Data:        hexutil.MustDecode(watchedEvent.Data),
		BlockNumber: uint64(watchedEvent.BlockNumber),
		TxHash:      common.HexToHash(watchedEvent.TxHash),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0x0"),
		Index:       uint(watchedEvent.Index),
		Removed:     false,
	}
}

func createTopics(topics ...string) []common.Hash {
	var topicsArray []common.Hash
	for _, topic := range topics {
		topicsArray = append(topicsArray, common.HexToHash(topic))
	}
	return topicsArray
}

func BigFromString(n string) *big.Int {
	b := new(big.Int)
	b.SetString(n, 10)
	return b
}

func GenerateSignature(s string) string {
	eventSignature := []byte(s)
	hash := crypto.Keccak256Hash(eventSignature)
	return hash.Hex()
}
