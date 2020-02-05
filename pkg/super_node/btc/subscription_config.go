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
	"math/big"

	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

/*
// HeaderModel is the db model for btc.header_cids table
// TxInput is the db model for btc.tx_inputs table
type TxInput struct {
	ID                    int64    `db:"id"`
	TxID                  int64    `db:"tx_id"`
	Index                 int64    `db:"index"`
	TxWitness             [][]byte `db:"tx_witness"`
	SignatureScript       []byte   `db:"sig_script"`
	PreviousOutPointHash  string   `db:"outpoint_hash"`
	PreviousOutPointIndex uint32   `db:"outpoint_index"`
}

// TxOutput is the db model for btc.tx_outputs table
type TxOutput struct {
	ID       int64  `db:"id"`
	TxID     int64  `db:"tx_id"`
	Index    int64  `db:"index"`
	Value    int64  `db:"value"`
	PkScript []byte `db:"pk_script"`
}

*/
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
	Off bool
	// Top level trx filters
	Index         int64    // allow filtering by index so that we can filter for only coinbase transactions (index 0) if we want to
	Segwit        bool     // allow filtering for segwit trxs
	WitnessHashes []string // allow filtering for specific witness hashes
	// TODO: trx input filters
	// TODO: trx output filters
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
	// Below default to false, which means we get all headers and no uncles by default
	sc.HeaderFilter = HeaderFilter{
		Off: viper.GetBool("superNode.btcSubscription.headerFilter.off"),
	}
	// Below defaults to false and two slices of length 0
	// Which means we get all transactions by default
	sc.TxFilter = TxFilter{
		Off:           viper.GetBool("superNode.btcSubscription.txFilter.off"),
		Index:         viper.GetInt64("superNode.btcSubscription.txFilter.index"),
		Segwit:        viper.GetBool("superNode.btcSubscription.txFilter.segwit"),
		WitnessHashes: viper.GetStringSlice("superNode.btcSubscription.txFilter.witnessHashes"),
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
