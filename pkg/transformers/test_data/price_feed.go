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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

var (
	medianizerAddress = common.HexToAddress("0x99041f808d598b782d5a3e498681c2452a31da08")
	blockNumber       = uint64(6147230)
	txIndex           = uint(119)
)

// https://etherscan.io/tx/0xa51a50a2adbfba4e2ab3d72dfd67a21c769f1bc8d2b180663a15500a56cde58f
var EthPriceFeedLog = types.Log{
	Address:     medianizerAddress,
	Topics:      []common.Hash{common.HexToHash("0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527")},
	Data:        common.FromHex("00000000000000000000000000000000000000000000001486f658319fb0c100"),
	BlockNumber: blockNumber,
	TxHash:      common.HexToHash("0xa51a50a2adbfba4e2ab3d72dfd67a21c769f1bc8d2b180663a15500a56cde58f"),
	TxIndex:     txIndex,
	BlockHash:   fakes.FakeHash,
	Index:       8,
	Removed:     false,
}

var rawPriceFeedLog, _ = json.Marshal(EthPriceFeedLog)
var PriceFeedModel = price_feeds.PriceFeedModel{
	BlockNumber:       blockNumber,
	MedianizerAddress: EthPriceFeedLog.Address.String(),
	UsdValue:          "378.659938889700015352",
	LogIndex:          EthPriceFeedLog.Index,
	TransactionIndex:  EthPriceFeedLog.TxIndex,
	Raw:               rawPriceFeedLog,
}
