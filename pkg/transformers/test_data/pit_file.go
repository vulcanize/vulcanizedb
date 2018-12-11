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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var EthPitFileDebtCeilingLog = types.Log{
	Address: common.HexToAddress(constants.PitContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x29ae811400000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x4c696e6500000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000001e240"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000004429ae81144c696e6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e240"),
	BlockNumber: 22,
	TxHash:      common.HexToHash("0xd744878a0b6655e3ba729e1019f56b563b4a16750196b8ad6104f3977db43f42"),
	TxIndex:     333,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var rawPitFileDebtCeilingLog, _ = json.Marshal(EthPitFileDebtCeilingLog)
var PitFileDebtCeilingModel = debt_ceiling.PitFileDebtCeilingModel{
	What:             "Line",
	Data:             shared.ConvertToWad(big.NewInt(123456).String()),
	LogIndex:         EthPitFileDebtCeilingLog.Index,
	TransactionIndex: EthPitFileDebtCeilingLog.TxIndex,
	Raw:              rawPitFileDebtCeilingLog,
}

var EthPitFileIlkLineLog = types.Log{
	Address: common.HexToAddress(constants.PitContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x1a0b287e00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6c696e6500000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b0000000000000000000000000000000000000000000000006c696e6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e8d4a51000"),
	BlockNumber: 12,
	TxHash:      common.HexToHash("0x2e27c962a697d4f7ec5d3206d0c058bd510f7593a711f082e55f3b62d44d8dd9"),
	TxIndex:     112,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var rawPitFileIlkLineLog, _ = json.Marshal(EthPitFileIlkLineLog)
var PitFileIlkLineModel = ilk2.PitFileIlkModel{
	Ilk:              "fake ilk",
	What:             "line",
	Data:             "0.000001000000000000",
	LogIndex:         EthPitFileIlkLineLog.Index,
	TransactionIndex: EthPitFileIlkLineLog.TxIndex,
	Raw:              rawPitFileIlkLineLog,
}

var EthPitFileIlkSpotLog = types.Log{
	Address: common.HexToAddress(constants.PitContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x1a0b287e00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x73706f7400000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b00000000000000000000000000000000000000000000000073706f7400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e8d4a51000"),
	BlockNumber: 11,
	TxHash:      common.HexToHash("0x1ba8125f60fa045c85b35df3983bee37db8627fbc32e3442a5cf17c85bb83f09"),
	TxIndex:     111,
	BlockHash:   fakes.FakeHash,
	Index:       14,
	Removed:     false,
}

var rawPitFileIlkSpotLog, _ = json.Marshal(EthPitFileIlkSpotLog)
var PitFileIlkSpotModel = ilk2.PitFileIlkModel{
	Ilk:              "fake ilk",
	What:             "spot",
	Data:             "0.000000000000001000000000000",
	LogIndex:         EthPitFileIlkSpotLog.Index,
	TransactionIndex: EthPitFileIlkSpotLog.TxIndex,
	Raw:              rawPitFileIlkSpotLog,
}
