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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"math/big"
)

var (
	TemporaryFrobAddress     = "0xff3f2400f1600f3f493a9a92704a29b96795af1a"
	TemporaryFrobBlockHash   = common.HexToHash("0x67ae45eace52de052a0fc58598974b101733f823fc191329ace7aded9a72b84b")
	TemporaryFrobBlockNumber = int64(13)
	TemporaryFrobData        = "0x000000000000000000000000000000000000000000000000000000000000000f0000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000019"
	TemporaryFrobTransaction = "0xbcff98316acb5732891d1a7e02f23ec12fbf8c231ca4b5530fa7a21c1e9b6aa9"
)

var (
	// need to set bytes as 0 or else the big Int 0 evaluates differently from the one unpacked by the abi
	art     = big.NewInt(20)
	dink    = big.NewInt(10)
	dart    = big.NewInt(0).SetBytes([]byte{0})
	iArt    = big.NewInt(25)
	frobLad = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 217, 34, 137, 65, 83, 190, 158, 239, 123, 114, 24, 220, 86, 93, 29, 12, 226, 160, 146}
	gem, _  = big.NewInt(0).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639926", 10)
	ink     = big.NewInt(15)
)

var EthFrobLog = types.Log{
	Address: common.HexToAddress(TemporaryFrobAddress),
	Topics: []common.Hash{
		common.HexToHash("0xb2afa28318bcc689926b52835d844de174ef8de97e982a85c0199d584920791b"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
	},
	Data:        hexutil.MustDecode(TemporaryFrobData),
	BlockNumber: uint64(TemporaryFrobBlockNumber),
	TxHash:      common.HexToHash(TemporaryFrobTransaction),
	TxIndex:     123,
	BlockHash:   TemporaryFrobBlockHash,
	Index:       0,
	Removed:     false,
}

var FrobEntity = frob.FrobEntity{
	Ilk:  ilk,
	Lad:  frobLad,
	Dink: dink,
	Dart: dart,
	Ink:  ink,
	Art:  art,
	IArt: iArt,
}

var FrobModel = frob.FrobModel{
	Ilk:  ilk[:],
	Lad:  frobLad[:],
	Dink: dink.String(),
	Dart: dart.String(),
	Ink:  ink.String(),
	Art:  art.String(),
	IArt: iArt.String(),
}
