// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package btc

import (
	"errors"
	"math/big"

	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// SubscriptionSettings config is used by a subscriber to specify what bitcoin data to stream from the super node
type SubscriptionSettings struct {
	BackFill     bool
	BackFillOnly bool
	Start        *big.Int
	End          *big.Int // set to 0 or a negative value to have no ending block
	HeaderFilter HeaderFilter
	TxFilter     TxFilter
}

// HeaderFilter contains filter settings for headers
type HeaderFilter struct {
	Off bool
}

// TxFilter contains filter settings for txs
type TxFilter struct {
	Off             bool
	Segwit          bool     // allow filtering for segwit trxs
	WitnessHashes   []string // allow filtering for specific witness hashes
	Indexes         []int64  // allow filtering for specific transaction indexes (e.g. 0 for coinbase transactions)
	PkScriptClasses []uint8  // allow filtering for txs that have at least one tx output with the specified pkscript class
	MultiSig        bool     // allow filtering for txs that have at least one tx output that requires more than one signature
	Addresses       []string // allow filtering for txs that have at least one tx output with at least one of the provided addresses
}

// Init is used to initialize a EthSubscription struct with env variables
func NewEthSubscriptionConfig() (*SubscriptionSettings, error) {
	sc := new(SubscriptionSettings)
	// Below default to false, which means we do not backfill by default
	sc.BackFill = viper.GetBool("superNode.btcSubscription.historicalData")
	sc.BackFillOnly = viper.GetBool("superNode.btcSubscription.historicalDataOnly")
	// Below default to 0
	// 0 start means we start at the beginning and 0 end means we continue indefinitely
	sc.Start = big.NewInt(viper.GetInt64("superNode.btcSubscription.startingBlock"))
	sc.End = big.NewInt(viper.GetInt64("superNode.btcSubscription.endingBlock"))
	// Below default to false, which means we get all headers by default
	sc.HeaderFilter = HeaderFilter{
		Off: viper.GetBool("superNode.btcSubscription.headerFilter.off"),
	}
	// Below defaults to false and two slices of length 0
	// Which means we get all transactions by default
	pksc := viper.Get("superNode.btcSubscription.txFilter.pkScriptClass")
	pkScriptClasses, ok := pksc.([]uint8)
	if !ok {
		return nil, errors.New("superNode.btcSubscription.txFilter.pkScriptClass needs to be an array of uint8s")
	}
	is := viper.Get("superNode.btcSubscription.txFilter.indexes")
	indexes, ok := is.([]int64)
	if !ok {
		return nil, errors.New("superNode.btcSubscription.txFilter.indexes needs to be an array of int64s")
	}
	sc.TxFilter = TxFilter{
		Off:             viper.GetBool("superNode.btcSubscription.txFilter.off"),
		Segwit:          viper.GetBool("superNode.btcSubscription.txFilter.segwit"),
		WitnessHashes:   viper.GetStringSlice("superNode.btcSubscription.txFilter.witnessHashes"),
		PkScriptClasses: pkScriptClasses,
		Indexes:         indexes,
		MultiSig:        viper.GetBool("superNode.btcSubscription.txFilter.multiSig"),
		Addresses:       viper.GetStringSlice("superNode.btcSubscription.txFilter.addresses"),
	}
	return sc, nil
}

// StartingBlock satisfies the SubscriptionSettings() interface
func (sc *SubscriptionSettings) StartingBlock() *big.Int {
	return sc.Start
}

// EndingBlock satisfies the SubscriptionSettings() interface
func (sc *SubscriptionSettings) EndingBlock() *big.Int {
	return sc.End
}

// HistoricalData satisfies the SubscriptionSettings() interface
func (sc *SubscriptionSettings) HistoricalData() bool {
	return sc.BackFill
}

// HistoricalDataOnly satisfies the SubscriptionSettings() interface
func (sc *SubscriptionSettings) HistoricalDataOnly() bool {
	return sc.BackFillOnly
}

// ChainType satisfies the SubscriptionSettings() interface
func (sc *SubscriptionSettings) ChainType() shared.ChainType {
	return shared.Bitcoin
}
