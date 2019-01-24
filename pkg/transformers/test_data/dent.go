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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"strconv"
)

var (
	DentData            = "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000645ff3a382000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000098a7d9b8314c000000000000000000000000000000000000000000000000000029a2241af62c0000"
	DentTransactionHash = "0x5a210319fcd31eea5959fedb4a1b20881c21a21976e23ff19dff3b44cc1c71e8"
	dentBidId           = int64(1)
	dentLot             = "11000000000000000000"
	dentBid             = "3000000000000000000"
	dentGuy             = "0x64d922894153BE9EEf7b7218dc565d1D0Ce2a092"
	dentRawJson, _      = json.Marshal(DentLog)
)

var DentLog = types.Log{
	Address: common.HexToAddress(KovanFlipperContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x5ff3a38200000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000098a7d9b8314c0000"),
	},
	Data:        hexutil.MustDecode(DentData),
	BlockNumber: 15,
	TxHash:      common.HexToHash(DentTransactionHash),
	TxIndex:     5,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var DentModel = dent.DentModel{
	BidId:            strconv.FormatInt(dentBidId, 10),
	Lot:              dentLot,
	Bid:              dentBid,
	Guy:              dentGuy,
	LogIndex:         DentLog.Index,
	TransactionIndex: DentLog.TxIndex,
	Raw:              dentRawJson,
}
