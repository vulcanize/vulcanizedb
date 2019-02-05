// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/repo"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"math/big"
)

var EthDripFileIlkLog = types.Log{
	Address: common.HexToAddress(KovanDripContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x1a0b287e00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b6520766f77000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b00000000000000000000000000000000000000000000000066616b6520766f77000000000000000000000000000000000000000000000000000000000000000000000000000000000000009B3F7188CE95D16E5AE0000000"),
	BlockNumber: 35,
	TxHash:      common.HexToHash("0xa1c31b7e6389470902237161263558615e60b40f2e63060b2f4aeafe92d57e5f"),
	TxIndex:     12,
	BlockHash:   fakes.FakeHash,
	Index:       15,
	Removed:     false,
}

var rawDripFileIlkLog, _ = json.Marshal(EthDripFileIlkLog)
var DripFileIlkModel = ilk2.DripFileIlkModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	Vow:              "66616b6520766f77000000000000000000000000000000000000000000000000",
	Tax:              "12300.000000000000000000000000000",
	LogIndex:         EthDripFileIlkLog.Index,
	TransactionIndex: EthDripFileIlkLog.TxIndex,
	Raw:              rawDripFileIlkLog,
}

var EthDripFileRepoLog = types.Log{
	Address: common.HexToAddress(KovanDripContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x29ae811400000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520776861740000000000000000000000000000000000000000000000"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000007b"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000004429ae811466616b6520776861740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b"),
	BlockNumber: 36,
	TxHash:      common.HexToHash("0xeeaa16de1d91c239b66773e8c2116a26cfeaaf5d962b31466c9bf047a5caa20f"),
	TxIndex:     13,
	BlockHash:   fakes.FakeHash,
	Index:       16,
	Removed:     false,
}

var rawDripFileRepoLog, _ = json.Marshal(EthDripFileRepoLog)
var DripFileRepoModel = repo.DripFileRepoModel{
	What:             "fake what",
	Data:             big.NewInt(123).String(),
	LogIndex:         EthDripFileRepoLog.Index,
	TransactionIndex: EthDripFileRepoLog.TxIndex,
	Raw:              rawDripFileRepoLog,
}

var EthDripFileVowLog = types.Log{
	Address: common.HexToAddress(KovanDripContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0xe9b674b900000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x766f770000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044e9b674b966616b652077686174000000000000000000000000000000000000000000000066616b6520646174610000000000000000000000000000000000000000000000"),
	BlockNumber: 51,
	TxHash:      common.HexToHash("0x586e26b71b41fcd6905044dbe8f0cca300517542278f74a9b925c4f800fed85c"),
	TxIndex:     14,
	BlockHash:   fakes.FakeHash,
	Index:       17,
	Removed:     false,
}

var rawDripFileVowLog, _ = json.Marshal(EthDripFileVowLog)
var DripFileVowModel = vow.DripFileVowModel{
	What:             "vow",
	Data:             "0x3728e9777B2a0a611ee0F89e00E01044ce4736d1",
	LogIndex:         EthDripFileVowLog.Index,
	TransactionIndex: EthDripFileVowLog.TxIndex,
	Raw:              rawDripFileVowLog,
}
