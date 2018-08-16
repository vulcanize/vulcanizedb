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
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
)

var tendLot = big.NewInt(100)
var tendBid = big.NewInt(50)
var tendGuy = common.HexToAddress("0x64d922894153be9eef7b7218dc565d1d0ce2a092")
var tic = new(big.Int).SetBytes([]byte{0})
var tendEra = big.NewInt(1533916180)
var RawJson, _ = json.Marshal(TendLog)
var rawString = string(RawJson)

var TendLog = types.Log{
	Address:     common.HexToAddress(FlipAddress),
	Topics:      []common.Hash{common.HexToHash(shared.TendSignature)},
	Data:        hexutil.MustDecode(TendData),
	BlockNumber: uint64(TendBlockNumber),
	TxHash:      common.HexToHash(TendTransactionHash),
	TxIndex:     1,
	BlockHash:   common.HexToHash(TendBlockHash),
	Index:       0,
	Removed:     false,
}

var tendId = int64(1)
var TendEntity = tend.TendEntity{
	Id:               big.NewInt(tendId),
	Lot:              tendLot,
	Bid:              tendBid,
	Guy:              tendGuy,
	Tic:              tic,
	Era:              tendEra,
	TransactionIndex: TendLog.TxIndex,
	Raw:              TendLog,
}

var TendModel = tend.TendModel{
	Id:               strconv.FormatInt(tendId, 10),
	Lot:              tendLot.String(),
	Bid:              tendBid.String(),
	Guy:              tendGuy[:],
	Tic:              tic.String(),
	Era:              time.Unix(tendEra.Int64(), 0),
	TransactionIndex: TendLog.TxIndex,
	Raw:              rawString,
}
