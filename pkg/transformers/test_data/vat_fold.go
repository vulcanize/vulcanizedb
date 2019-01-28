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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_fold"
)

var EthVatFoldLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
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
	Ilk:              "5245500000000000000000000000000000000000000000000000000000000000",
	Urn:              "0x3728e9777B2a0a611ee0F89e00E01044ce4736d1",
	Rate:             "0.000000000000000000000000002",
	LogIndex:         EthVatFoldLog.Index,
	TransactionIndex: EthVatFoldLog.TxIndex,
	Raw:              rawVatFoldLog,
}
