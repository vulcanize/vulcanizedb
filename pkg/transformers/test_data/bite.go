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
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"math/big"
)

const (
	TemporaryBiteBlockNumber = int64(26)
	TemporaryBiteData        = "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005"
	TemporaryBiteTransaction = "0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"
)

var (
	biteInk        = big.NewInt(1)
	biteArt        = big.NewInt(2)
	biteTab        = big.NewInt(3)
	biteFlip       = big.NewInt(4)
	biteIArt       = big.NewInt(5)
	biteRawJson, _ = json.Marshal(EthBiteLog)
	biteIlk        = [32]byte{69, 84, 72, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	biteLad        = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 216, 180, 20, 126, 218, 128, 254, 199, 18, 42, 225, 109, 162, 71, 156, 189, 127, 251}
	biteIlkString  = "ETH"
	biteLadString  = "0x0000d8b4147eDa80Fec7122AE16DA2479Cbd7ffB"
)

var EthBiteLog = types.Log{
	Address: common.HexToAddress("0x2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
	Topics: []common.Hash{
		common.HexToHash("0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"),
	},
	Data:        hexutil.MustDecode(TemporaryBiteData),
	BlockNumber: uint64(TemporaryBiteBlockNumber),
	TxHash:      common.HexToHash(TemporaryBiteTransaction),
	TxIndex:     111,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}

var BiteEntity = bite.BiteEntity{
	Ilk:              biteIlk,
	Urn:              biteLad,
	Ink:              biteInk,
	Art:              biteArt,
	Tab:              biteTab,
	Flip:             biteFlip,
	IArt:             biteIArt,
	LogIndex:         EthBiteLog.Index,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              EthBiteLog,
}

var BiteModel = bite.BiteModel{
	Ilk:              biteIlkString,
	Urn:              biteLadString,
	Ink:              biteInk.String(),
	Art:              biteArt.String(),
	IArt:             biteIArt.String(),
	Tab:              biteTab.String(),
	NFlip:            biteFlip.String(),
	LogIndex:         EthBiteLog.Index,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              biteRawJson,
}
