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

package price_feeds

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogValueEntity struct {
	Val common.Address
}

type PriceFeedModel struct {
	BlockNumber       uint64 `db:"block_number"`
	MedianizerAddress string `db:"medianizer_address"`
	UsdValue          string `db:"usd_value"`
	LogIndex          uint   `db:"log_idx"`
	TransactionIndex  uint   `db:"tx_idx"`
	Raw               []byte `db:"raw_log"`
}

func Convert(conversion string, value string, prec int) string {
	var bgflt = big.NewFloat(0.0)
	bgflt.SetString(value)
	switch conversion {
	case "ray":
		bgflt.Quo(bgflt, Ray)
	case "wad":
		bgflt.Quo(bgflt, Ether)
	}
	return bgflt.Text('g', prec)
}
