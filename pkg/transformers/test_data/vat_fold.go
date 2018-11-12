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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
)

var EthVatFoldLog = types.Log{
	Address: common.HexToAddress(constants.VatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0xe6a6a64d00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000064e6a6a64d45544800000000000000000000000000000000000000000000000000000000000000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d10000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 8940380,
	TxHash:      common.HexToHash("0xfb37b7a88aa8ad14538d1e244a55939fa07c1828e5ca8168bf4edd56f5fc4d57"),
	TxIndex:     8,
	BlockHash:   fakes.FakeHash,
	Index:       5,
	Removed:     false,
}

var rawVatFoldLog, _ = json.Marshal(EthVatFoldLog)
var VatFoldModel = vat_fold.VatFoldModel{
	Ilk:              "REP",
	Urn:              "0x3728e9777B2a0a611ee0F89e00E01044ce4736d1",
	Rate:             "2",
	LogIndex:         EthVatFoldLog.Index,
	TransactionIndex: EthVatFoldLog.TxIndex,
	Raw:              rawVatFoldLog,
}
