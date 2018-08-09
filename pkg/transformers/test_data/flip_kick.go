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
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"time"
)

var (
	TemporaryFlipBlockNumber     = int64(10)
	TemporaryFlipAddress         = "0x08cb6176addcca2e1d1ffe21bee464b72ee4cd8d"
	TemporaryFlipKickBlockHash   = "0x32f8b12023b3a1b4c73f9a46da976931b0355714ada8b8044ebcb2cd295751a9"
	TemporaryFlipKickData        = "0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000008cb6176addcca2e1d1ffe21bee464b72ee4cd8d00000000000000000000000038219779a699d67d7e7740b8c8f43d3e2dae218266616b6520696c6b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000000000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a09200000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977000000000000000000000000000000000000000000000000000000005b69b8e7000000000000000000000000000000000000000000000000000000005b607e670000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca0000000000000000000000000000000000000000000000000000000000000032"
	TemporaryFlipKickTransaction = "0x6b155a55fd77b751195deeebf7abfd8691ca01ee588817a920f19d5b27f65191"
)

var idString = "1"
var id, _ = new(big.Int).SetString(idString, 10)
var mom = "0x08cb6176addcca2e1d1ffe21bee464b72ee4cd8d"
var vat = "0x38219779a699d67d7e7740b8c8f43d3e2dae2182"
var ilk = [32]byte{102, 97, 107, 101, 32, 105, 108, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var lotString = "100"
var lot, _ = new(big.Int).SetString(lotString, 10)
var bidString = "0"
var bid = new(big.Int).SetBytes([]byte{0})
var guy = "0x64d922894153be9eef7b7218dc565d1d0ce2a092"
var gal = "0x07fa9ef6609ca7921112231f8f195138ebba2977"
var end = int64(1533655271)
var era = int64(1533050471)
var lad = "0x7340e006f4135ba6970d43bf43d88dcad4e7a8ca"
var tabString = "50"
var tab, _ = new(big.Int).SetString(tabString, 10)

var EthFlipKickLog = types.Log{
	Address:     common.HexToAddress(TemporaryFlipAddress),
	Topics:      []common.Hash{common.HexToHash(flip_kick.FlipKickSignature)},
	Data:        hexutil.MustDecode(TemporaryFlipKickData),
	BlockNumber: uint64(TemporaryFlipBlockNumber),
	TxHash:      common.HexToHash(TemporaryFlipKickTransaction),
	TxIndex:     0,
	BlockHash:   common.HexToHash(TemporaryFlipKickBlockHash),
	Index:       0,
	Removed:     false,
}

var FlipKickEntity = flip_kick.FlipKickEntity{
	Id:  id,
	Mom: common.HexToAddress(mom),
	Vat: common.HexToAddress(vat),
	Ilk: ilk,
	Lot: lot,
	Bid: bid,
	Guy: common.HexToAddress(guy),
	Gal: common.HexToAddress(gal),
	End: big.NewInt(end),
	Era: big.NewInt(era),
	Lad: common.HexToAddress(lad),
	Tab: tab,
}

var FlipKickModel = flip_kick.FlipKickModel{
	Id:  idString,
	Mom: mom,
	Vat: vat,
	Ilk: "0x" + common.Bytes2Hex(ilk[:]),
	Lot: lotString,
	Bid: bidString,
	Guy: guy,
	Gal: gal,
	End: time.Unix(end, 0),
	Era: time.Unix(era, 0),
	Lad: lad,
	Tab: tabString,
}

type FlipKickDBRow struct {
	DbID     int64 `db:"db_id"`
	HeaderId int64 `db:"header_id"`
	flip_kick.FlipKickModel
}
