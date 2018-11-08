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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
)

var (
	tendBidId           = int64(10)
	tendLot             = "85000000000000000000"
	tendBid             = "1000000000000000000"
	tendGuy             = "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6"
	tendData            = "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000644b43ed12000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000049b9ca9a6943400000000000000000000000000000000000000000000000000000de0b6b3a7640000"
	tendTransactionHash = "0x7909c8793ded2b8348f5db623044fbc26bb7ab78ad5792897abdf68ddc1df63d"
	tendBlockHash       = "0xa8ea87147c0a68daeb6b1d9f8c0937ba975a650809cab80d19c969e8d0df452c"
	TendTic             = "0"
)

var TendLogNote = types.Log{
	Address: common.HexToAddress(constants.FlipperContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x4b43ed1200000000000000000000000000000000000000000000000000000000"), //abbreviated tend function signature
		common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"), //msg caller address
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"), //first param of the function called (i.e. flip kick id)
		common.HexToHash("0x0000000000000000000000000000000000000000000000049b9ca9a694340000"), //second param of the function called (i.e. lot)
	},
	Data:        hexutil.MustDecode(tendData),
	BlockNumber: 11,
	TxHash:      common.HexToHash(tendTransactionHash),
	TxIndex:     10,
	BlockHash:   common.HexToHash(tendBlockHash),
	Index:       1,
	Removed:     false,
}
var RawLogNoteJson, _ = json.Marshal(TendLogNote)

var TendModel = tend.TendModel{
	BidId:            strconv.FormatInt(tendBidId, 10),
	Lot:              tendLot,
	Bid:              tendBid,
	Guy:              tendGuy,
	Tic:              TendTic,
	LogIndex:         TendLogNote.Index,
	TransactionIndex: TendLogNote.TxIndex,
	Raw:              string(RawLogNoteJson),
}
