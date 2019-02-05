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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
)

var EthVatSlipLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
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
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var rawVatSlipLog, _ = json.Marshal(EthVatSlipLog)
var VatSlipModel = vat_slip.VatSlipModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	Guy:              "0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6",
	Rad:              "987654321",
	TransactionIndex: EthVatSlipLog.TxIndex,
	LogIndex:         EthVatSlipLog.Index,
	Raw:              rawVatSlipLog,
}
