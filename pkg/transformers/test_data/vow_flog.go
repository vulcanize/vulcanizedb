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
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vow_flog"
)

var EthVowFlogLog = types.Log{
	Address: common.HexToAddress(KovanVowContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x35aee16f00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000008e84a1e068d77059cbe263c43ad0cdc130863313"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000539"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000002435aee16f0000000000000000000000000000000000000000000000000000000000000539"),
	BlockNumber: 11,
	TxHash:      common.HexToHash("0x47ffd75c1cda1d5c08219755743663a3790e4f5ae9e1fcb85bb3fe0d74bb7109"),
	TxIndex:     4,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}

var rawVowFlogLog, _ = json.Marshal(EthVowFlogLog)
var VowFlogModel = vow_flog.VowFlogModel{
	Era:              "1337",
	LogIndex:         EthVowFlogLog.Index,
	TransactionIndex: EthVowFlogLog.TxIndex,
	Raw:              rawVowFlogLog,
}
