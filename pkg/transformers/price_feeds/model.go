// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
