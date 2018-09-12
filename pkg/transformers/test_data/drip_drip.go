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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_drip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var EthDripDripLog = types.Log{
	Address: common.HexToAddress(shared.DripContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x44e2a5a800000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000002444e2a5a866616b6520696c6b000000000000000000000000000000000000000000000000"),
	BlockNumber: 62,
	TxHash:      common.HexToHash("0xa34fd5cfcb125ebfc81d33633495701b531753669712092bdb8aa6159a240040"),
	TxIndex:     10,
	BlockHash:   common.HexToHash("0x7a2aa72468986774d90dc8c80d436956a5b6a7e5acea430fb8a79f9217ef00a3"),
	Index:       0,
	Removed:     false,
}

var rawDripDripLog, _ = json.Marshal(EthDripDripLog)
var DripDripModel = drip_drip.DripDripModel{
	Ilk:              "fake ilk",
	TransactionIndex: EthDripDripLog.TxIndex,
	Raw:              rawDripDripLog,
}
