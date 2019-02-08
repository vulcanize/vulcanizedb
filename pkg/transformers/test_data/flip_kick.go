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
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
)

var (
	idString  = "1"
	id, _     = new(big.Int).SetString(idString, 10)
	lotString = "100"
	lot, _    = new(big.Int).SetString(lotString, 10)
	bidString = "0"
	bid       = new(big.Int).SetBytes([]byte{0})
	gal       = "0x07Fa9eF6609cA7921112231F8f195138ebbA2977"
	end       = int64(1535991025)
	urnBytes  = [32]byte{115, 64, 224, 6, 244, 19, 91, 166, 151, 13, 67, 191, 67, 216, 141, 202, 212, 231, 168, 202, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	urnString = "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA"
	tabString = "50"
	tab, _    = new(big.Int).SetString(tabString, 10)
	rawLog, _ = json.Marshal(EthFlipKickLog)
)

var (
	flipKickTransactionHash = "0xd11ab35cfb1ad71f790d3dd488cc1a2046080e765b150e8997aa0200947d4a9b"
	flipKickData            = "0x0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977000000000000000000000000000000000000000000000000000000005b8d5cf10000000000000000000000000000000000000000000000000000000000000032"
	FlipKickBlockNumber     = int64(10)
)

var EthFlipKickLog = types.Log{
	Address: common.HexToAddress(KovanFlipperContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0xbac86238bdba81d21995024470425ecb370078fa62b7271b90cf28cbd1e3e87e"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		common.HexToHash("0x7340e006f4135ba6970d43bf43d88dcad4e7a8ca000000000000000000000000"),
	},
	Data:        hexutil.MustDecode(flipKickData),
	BlockNumber: uint64(FlipKickBlockNumber),
	TxHash:      common.HexToHash(flipKickTransactionHash),
	TxIndex:     999,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var FlipKickEntity = flip_kick.FlipKickEntity{
	Id:               id,
	Lot:              lot,
	Bid:              bid,
	Gal:              common.HexToAddress(gal),
	End:              big.NewInt(end),
	Urn:              urnBytes,
	Tab:              tab,
	TransactionIndex: EthFlipKickLog.TxIndex,
	LogIndex:         EthFlipKickLog.Index,
	Raw:              EthFlipKickLog,
}

var FlipKickModel = flip_kick.FlipKickModel{
	BidId:            idString,
	Lot:              lotString,
	Bid:              bidString,
	Gal:              gal,
	End:              time.Unix(end, 0),
	Urn:              urnString,
	Tab:              tabString,
	TransactionIndex: EthFlipKickLog.TxIndex,
	LogIndex:         EthFlipKickLog.Index,
	Raw:              rawLog,
}

type FlipKickDBRow struct {
	ID       int64
	HeaderId int64 `db:"header_id"`
	flip_kick.FlipKickModel
}
