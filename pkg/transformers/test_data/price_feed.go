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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var (
	medianizerAddress = common.HexToAddress("0x99041f808d598b782d5a3e498681c2452a31da08")
	blockNumber       = uint64(6147230)
	txIndex           = uint(119)
)

// https://etherscan.io/tx/0xa51a50a2adbfba4e2ab3d72dfd67a21c769f1bc8d2b180663a15500a56cde58f
var EthPriceFeedLog = types.Log{
	Address:     medianizerAddress,
	Topics:      []common.Hash{common.HexToHash(shared.LogValueSignature)},
	Data:        common.FromHex("00000000000000000000000000000000000000000000001486f658319fb0c100"),
	BlockNumber: blockNumber,
	TxHash:      common.HexToHash("0xa51a50a2adbfba4e2ab3d72dfd67a21c769f1bc8d2b180663a15500a56cde58f"),
	TxIndex:     txIndex,
	BlockHash:   common.HexToHash("0x27ecebbf69eefa3bb3cf65f472322a80ff4946653a50a2171dc605f49829467d"),
	Index:       0,
	Removed:     false,
}

var rawPriceFeedLog, _ = json.Marshal(EthPriceFeedLog)
var PriceFeedModel = price_feeds.PriceFeedModel{
	BlockNumber:       blockNumber,
	MedianizerAddress: EthPriceFeedLog.Address[:],
	UsdValue:          "378.6599388897",
	TransactionIndex:  EthPriceFeedLog.TxIndex,
}
