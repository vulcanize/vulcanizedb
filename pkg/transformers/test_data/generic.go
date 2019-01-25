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

package test_data

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
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

var GenericTestConfig = shared_t.TransformerConfig{
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
