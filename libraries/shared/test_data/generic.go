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

package test_data

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

type GenericModel struct{}
type GenericEntity struct{}

var startingBlockNumber = rand.Int63()
var topic0 = "0x" + randomString(64)

var GenericTestLog = func() types.Log {
	return types.Log{
		Address:     fakeAddress(),
		Topics:      []common.Hash{common.HexToHash(topic0), fakeHash()},
		Data:        hexutil.MustDecode(fakeHash().Hex()),
		BlockNumber: uint64(startingBlockNumber),
		TxHash:      fakeHash(),
		TxIndex:     uint(rand.Int31()),
		BlockHash:   fakeHash(),
		Index:       uint(rand.Int31()),
	}
}

var GenericTestConfig = transformer.EventTransformerConfig{
	TransformerName:     "generic-test-transformer",
	ContractAddresses:   []string{fakeAddress().Hex()},
	ContractAbi:         randomString(100),
	Topic:               topic0,
	StartingBlockNumber: startingBlockNumber,
	EndingBlockNumber:   startingBlockNumber + 1,
}

func fakeAddress() common.Address {
	return common.HexToAddress("0x" + randomString(40))
}

func fakeHash() common.Hash {
	return common.HexToHash("0x" + randomString(64))
}

func randomString(length int) string {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	charset := "abcdef1234567890"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
