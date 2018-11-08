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
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/repo"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"math/big"
)

var EthDripFileIlkLog = types.Log{
	Address: common.HexToAddress(constants.DripContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x1a0b287e00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b6520766f77000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b00000000000000000000000000000000000000000000000066616b6520766f77000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b"),
	BlockNumber: 35,
	TxHash:      common.HexToHash("0xa1c31b7e6389470902237161263558615e60b40f2e63060b2f4aeafe92d57e5f"),
	TxIndex:     12,
	BlockHash:   common.HexToHash("0x0188f3ee3cc05aa72457fa328e6a461de31e4cbd429fc37f9a52da4e9773c0b4"),
	Index:       15,
	Removed:     false,
}

var rawDripFileIlkLog, _ = json.Marshal(EthDripFileIlkLog)
var DripFileIlkModel = ilk2.DripFileIlkModel{
	Ilk:              "fake ilk",
	Vow:              "fake vow",
	Tax:              big.NewInt(123).String(),
	LogIndex:         EthDripFileIlkLog.Index,
	TransactionIndex: EthDripFileIlkLog.TxIndex,
	Raw:              rawDripFileIlkLog,
}

var EthDripFileRepoLog = types.Log{
	Address: common.HexToAddress(constants.DripContractAddress),
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
	BlockHash:   common.HexToHash("0x89de4145ea8e34dfd9db9a7ea34f5be6f1f402e812fd389acca342513b353288"),
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
	Address: common.HexToAddress(constants.DripContractAddress),
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
	BlockHash:   common.HexToHash("0xbec69b1e93503679c9c006819477b86fe16aaff3a418da1e916c431b68be5522"),
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
