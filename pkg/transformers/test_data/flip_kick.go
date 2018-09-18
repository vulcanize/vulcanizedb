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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var (
	idString      = "1"
	id, _         = new(big.Int).SetString(idString, 10)
	lotString     = "100"
	lot, _        = new(big.Int).SetString(lotString, 10)
	bidString     = "0"
	bid           = new(big.Int).SetBytes([]byte{0})
	gal           = "0x07Fa9eF6609cA7921112231F8f195138ebbA2977"
	end           = int64(1535991025)
	urn           = [32]byte{115, 64, 224, 6, 244, 19, 91, 166, 151, 13, 67, 191, 67, 216, 141, 202, 212, 231, 168, 202, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	urnString     = "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA"
	tabString     = "50"
	tab, _        = new(big.Int).SetString(tabString, 10)
	rawLogJson, _ = json.Marshal(EthFlipKickLog)
	rawLogString  = string(rawLogJson)
)

var (
	flipKickTransactionHash = "0xd11ab35cfb1ad71f790d3dd488cc1a2046080e765b150e8997aa0200947d4a9b"
	flipKickData            = "0x0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977000000000000000000000000000000000000000000000000000000005b8d5cf10000000000000000000000000000000000000000000000000000000000000032"
	flipKickBlockHash       = "0x40fcad7863ab4bef421d638b7ad6116e81577f14a62ef847b07f8527944466fd"
	FlipKickBlockNumber     = int64(10)
)

var EthFlipKickLog = types.Log{
	Address: common.HexToAddress(shared.FlipperContractAddress),
	Topics: []common.Hash{
		common.HexToHash(shared.FlipKickSignature),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0x7340e006f4135ba6970d43bf43d88dcad4e7a8ca000000000000000000000000"),
	},
	Data:        hexutil.MustDecode(flipKickData),
	BlockNumber: uint64(FlipKickBlockNumber),
	TxHash:      common.HexToHash(flipKickTransactionHash),
	TxIndex:     999,
	BlockHash:   common.HexToHash(flipKickBlockHash),
	Index:       0,
	Removed:     false,
}

var FlipKickEntity = flip_kick.FlipKickEntity{
	Id:               id,
	Lot:              lot,
	Bid:              bid,
	Gal:              common.HexToAddress(gal),
	End:              big.NewInt(end),
	Urn:              urn,
	Tab:              tab,
	TransactionIndex: EthFlipKickLog.TxIndex,
	Raw:              EthFlipKickLog,
}

var FlipKickModel = flip_kick.FlipKickModel{
	BidId:            idString,
	Lot:              lotString,
	Bid:              bidString,
	Gal:              gal,
	End:              time.Unix(end, 0),
	Urn:              urnString,
	Tab:              tabString,
	TransactionIndex: EthFlipKickLog.TxIndex,
	Raw:              rawLogString,
}

type FlipKickDBRow struct {
	ID       int64
	HeaderId int64 `db:"header_id"`
	flip_kick.FlipKickModel
}
