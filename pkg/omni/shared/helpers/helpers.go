// VulcanizeDB
// Copyright © 2018 Vulcanize

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
