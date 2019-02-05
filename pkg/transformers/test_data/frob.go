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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"math/big"
)

var (
	TemporaryFrobBlockNumber = int64(13)
	TemporaryFrobData        = "0x000000000000000000000000000000000000000000000000000000000000000f0000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000019"
	TemporaryFrobTransaction = "0xbcff98316acb5732891d1a7e02f23ec12fbf8c231ca4b5530fa7a21c1e9b6aa9"
)

var (
	// need to set bytes as 0 or else the big Int 0 evaluates differently from the one unpacked by the abi
	art           = big.NewInt(20)
	dink          = big.NewInt(10)
	dart          = big.NewInt(0).SetBytes([]byte{0})
	iArt          = big.NewInt(25)
	frobLad       = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 217, 34, 137, 65, 83, 190, 158, 239, 123, 114, 24, 220, 86, 93, 29, 12, 226, 160, 146}
	ink           = big.NewInt(15)
	ilk           = [32]byte{102, 97, 107, 101, 32, 105, 108, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	frobIlkString = "66616b6520696c6b000000000000000000000000000000000000000000000000"
	frobUrnString = "00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"
)

var EthFrobLog = types.Log{
	Address: common.HexToAddress(KovanPitContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0xb2afa28318bcc689926b52835d844de174ef8de97e982a85c0199d584920791b"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
	},
	Data:        hexutil.MustDecode(TemporaryFrobData),
	BlockNumber: uint64(TemporaryFrobBlockNumber),
	TxHash:      common.HexToHash(TemporaryFrobTransaction),
	TxIndex:     123,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}

var FrobEntity = frob.FrobEntity{
	Ilk:              ilk,
	Urn:              frobLad,
	Ink:              ink,
	Art:              art,
	Dink:             dink,
	Dart:             dart,
	IArt:             iArt,
	LogIndex:         EthFrobLog.Index,
	TransactionIndex: EthFrobLog.TxIndex,
	Raw:              EthFrobLog,
}

var rawFrobLog, _ = json.Marshal(EthFrobLog)
var FrobModel = frob.FrobModel{
	Ilk:              frobIlkString,
	Urn:              frobUrnString,
	Ink:              ink.String(),
	Art:              art.String(),
	Dink:             dink.String(),
	Dart:             dart.String(),
	IArt:             iArt.String(),
	LogIndex:         EthFrobLog.Index,
	TransactionIndex: EthFrobLog.TxIndex,
	Raw:              rawFrobLog,
}
