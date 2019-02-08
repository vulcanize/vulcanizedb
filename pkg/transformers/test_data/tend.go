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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
)

var (
	tendBidId           = int64(10)
	tendLot             = "85000000000000000000"
	tendBid             = "1000000000000000000"
	tendGuy             = "0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"
	tendData            = "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000644b43ed12000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000049b9ca9a6943400000000000000000000000000000000000000000000000000000de0b6b3a7640000"
	tendTransactionHash = "0x7909c8793ded2b8348f5db623044fbc26bb7ab78ad5792897abdf68ddc1df63d"
	tendBlockHash       = "0xa8ea87147c0a68daeb6b1d9f8c0937ba975a650809cab80d19c969e8d0df452c"
)

var TendLogNote = types.Log{
	Address: common.HexToAddress(KovanFlipperContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x4b43ed1200000000000000000000000000000000000000000000000000000000"), //abbreviated tend function signature
		common.HexToHash("0x0000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6"), //msg caller address
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"), //first param of the function called (i.e. flip kick id)
		common.HexToHash("0x0000000000000000000000000000000000000000000000049b9ca9a694340000"), //second param of the function called (i.e. lot)
	},
	Data:        hexutil.MustDecode(tendData),
	BlockNumber: 11,
	TxHash:      common.HexToHash(tendTransactionHash),
	TxIndex:     10,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var rawTendLog, _ = json.Marshal(TendLogNote)
var TendModel = tend.TendModel{
	BidId:            strconv.FormatInt(tendBidId, 10),
	Lot:              tendLot,
	Bid:              tendBid,
	Guy:              tendGuy,
	LogIndex:         TendLogNote.Index,
	TransactionIndex: TendLogNote.TxIndex,
	Raw:              rawTendLog,
}
