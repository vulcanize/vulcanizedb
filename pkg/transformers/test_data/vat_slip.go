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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
)

var EthVatSlipLog = types.Log{
	Address: common.HexToAddress(shared.VatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x42066cbb00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000003ade68b1"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000006442066cbb66616b6520696c6b0000000000000000000000000000000000000000000000000000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6000000000000000000000000000000000000000000000000000000003ade68b1"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0xb114ba306c80c86d51bdbf4a5ac8ed151020cd81b70cfa1dc9822f4a1f73930b"),
	TxIndex:     3,
	BlockHash:   common.HexToHash("0x34b7e5ddb3be73257a5a0087f10b8bf68d4df5c8831ec04c63ecae4094de72ad"),
	Index:       0,
	Removed:     false,
}

var rawVatSlipLog, _ = json.Marshal(EthVatSlipLog)
var VatSlipModel = vat_slip.VatSlipModel{
	Ilk:              "fake ilk",
	Guy:              common.HexToAddress("0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6").String(),
	Rad:              "987654321",
	TransactionIndex: EthVatSlipLog.TxIndex,
	Raw:              rawVatSlipLog,
}
