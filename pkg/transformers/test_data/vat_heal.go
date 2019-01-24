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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_heal"
)

var VatHealLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x990a5f6300000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000078"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000064990a5f637d7bee5fcfd8028cf7b00876c5b1421c800561a600000000000000000000000074686520762076616c75650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000078"),
	BlockNumber: 10,
	TxHash:      common.HexToHash("0x991b8079b1333024000dcaf2b00c24c5db0315e112a4ac4d912aa96a602e12b9"),
	TxIndex:     2,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}

var rawHealLog, _ = json.Marshal(VatHealLog)
var VatHealModel = vat_heal.VatHealModel{
	Urn:              "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6",
	V:                "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	Rad:              "120",
	LogIndex:         VatHealLog.Index,
	TransactionIndex: VatHealLog.TxIndex,
	Raw:              rawHealLog,
}
