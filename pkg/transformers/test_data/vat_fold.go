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
	"bytes"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
)

var EthVatFoldLog = types.Log{
	Address: common.HexToAddress(shared.VatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0xe6a6a64d00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000064e6a6a64d45544800000000000000000000000000000000000000000000000000000000000000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d10000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 72,
	TxHash:      common.HexToHash("0xe8f39fbb7fea3621f543868f19b1114e305aff6a063a30d32835ff1012526f91"),
	TxIndex:     8,
	BlockHash:   common.HexToHash("0xe3dd2e05bd8b92833e20ed83e2171bbc06a9ec823232eca1730a807bd8f5edc0"),
	Index:       5,
	Removed:     false,
}

var rawVatFoldLog, _ = json.Marshal(EthVatFoldLog)
var VatFoldModel = vat_fold.VatFoldModel{
	Ilk:              string(bytes.Trim(EthVatFoldLog.Topics[1].Bytes(), "\x00")),
	Urn:              string(bytes.Trim(EthVatFoldLog.Topics[2].Bytes(), "\x00")),
	Rate:             string(bytes.Trim(EthVatFoldLog.Topics[3].Bytes(), "\x00")),
	TransactionIndex: EthVatFoldLog.TxIndex,
	Raw:              rawVatFoldLog,
}
