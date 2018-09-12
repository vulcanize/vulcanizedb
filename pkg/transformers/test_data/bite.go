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
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"math/big"
	"strconv"
)

var (
	TemporaryBiteBlockHash   = common.HexToHash("0xd130caaccc9203ca63eb149faeb013aed21f0317ce23489c0486da2f9adcd0eb")
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
	biteRawString  = string(biteRawJson)
	biteIlk        = [32]byte{102, 97, 107, 101, 32, 105, 108, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	biteLad        = [32]byte{102, 97, 107, 101, 32, 108, 97, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	biteId         = int64(1)
)

var EthBiteLog = types.Log{
	Address: common.HexToAddress(shared.CatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b65206c6164000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode(TemporaryBiteData),
	BlockNumber: uint64(TemporaryBiteBlockNumber),
	TxHash:      common.HexToHash(TemporaryBiteTransaction),
	TxIndex:     111,
	BlockHash:   TemporaryBiteBlockHash,
	Index:       0,
	Removed:     false,
}

var BiteEntity = bite.BiteEntity{
	Id:               big.NewInt(biteId),
	Ilk:              biteIlk,
	Urn:              biteLad,
	Ink:              biteInk,
	Art:              biteArt,
	Tab:              biteTab,
	Flip:             biteFlip,
	IArt:             biteIArt,
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              EthBiteLog,
}

var BiteModel = bite.BiteModel{
	Id:               strconv.FormatInt(biteId, 10),
	Ilk:              biteIlk[:],
	Urn:              biteLad[:],
	Ink:              biteInk.String(),
	Art:              biteArt.String(),
	Tab:              biteTab.String(),
	Flip:             biteFlip.String(),
	IArt:             biteIArt.String(),
	TransactionIndex: EthBiteLog.TxIndex,
	Raw:              biteRawString,
}
