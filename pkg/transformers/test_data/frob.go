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
	TemporaryFrobBlockHash   = common.HexToHash("0xe1c4264e245ac31d4aed678df007199cffcfcaea7a75aeecc45122957abf4298")
	TemporaryFrobBlockNumber = int64(12)
	TemporaryFrobData        = "0x66616b6520696c6b00000000000000000000000000000000000000000000000000000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005b6b39e4"
	TemporaryFrobTransaction = "0xbcff98316acb5732891d1a7e02f23ec12fbf8c231ca4b5530fa7a21c1e9b6aa9"
)

var (
	// need to set bytes as 0 or else the big Int 0 evaluates differently from the one unpacked by the abi
	art     = big.NewInt(0).SetBytes([]byte{0})
	frobEra = big.NewInt(1533753828)
	frobLad = common.HexToAddress("0x64d922894153BE9EEf7b7218dc565d1D0Ce2a092")
	gem, _  = big.NewInt(0).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639926", 10)
	ink     = big.NewInt(10)
)

var EthFrobLog = types.Log{
	Address:     common.HexToAddress(TemporaryFrobAddress),
	Topics:      []common.Hash{common.HexToHash(frob.FrobEventSignature)},
	Data:        hexutil.MustDecode(TemporaryFrobData),
	BlockNumber: uint64(TemporaryFrobBlockNumber),
	TxHash:      common.HexToHash(TemporaryFrobTransaction),
	TxIndex:     123,
	BlockHash:   TemporaryFrobBlockHash,
	Index:       1,
	Removed:     false,
}

var FrobEntity = frob.FrobEntity{
	Ilk: ilk,
	Lad: frobLad,
	Gem: gem,
	Ink: ink,
	Art: art,
	Era: frobEra,
}

var FrobModel = frob.FrobModel{
	Ilk: ilk[:],
	Lad: frobLad[:],
	Gem: gem.String(),
	Ink: ink.String(),
	Art: art.String(),
	Era: frobEra.String(),
}
