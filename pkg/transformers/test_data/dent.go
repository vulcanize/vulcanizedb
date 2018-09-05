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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"strconv"
)

var (
	DentData            = "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000645ff3a382000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000098a7d9b8314c000000000000000000000000000000000000000000000000000029a2241af62c0000"
	DentTransactionHash = "0x5a210319fcd31eea5959fedb4a1b20881c21a21976e23ff19dff3b44cc1c71e8"
	DentBlockHash       = "0x105b771e04d7b8516f9291b1f006c46c09cfbff9efa8bc52498b171ff99d28b5"
	dentBidId           = int64(1)
	dentLot             = "11000000000000000000"
	dentBid             = "3000000000000000000"
	DentTic             = "0"
	dentGuy             = "0x64d922894153BE9EEf7b7218dc565d1D0Ce2a092"
	dentRawJson, _      = json.Marshal(DentLog)
)

var DentLog = types.Log{
	Address: common.HexToAddress(shared.FlipperContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x5ff3a38200000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000098a7d9b8314c0000"),
	},
	Data:        hexutil.MustDecode(DentData),
	BlockNumber: 15,
	TxHash:      common.HexToHash(DentTransactionHash),
	TxIndex:     5,
	BlockHash:   common.HexToHash(DentBlockHash),
	Index:       2,
	Removed:     false,
}

var DentModel = dent.DentModel{
	BidId:            strconv.FormatInt(dentBidId, 10),
	Lot:              dentLot,
	Bid:              dentBid,
	Guy:              dentGuy,
	Tic:              DentTic,
	TransactionIndex: DentLog.TxIndex,
	Raw:              dentRawJson,
}
