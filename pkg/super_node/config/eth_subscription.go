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

package config

import (
	"math/big"

	"github.com/spf13/viper"
)

// EthSubscription config is used by a subscriber to specify what eth data to stream from the super node
type EthSubscription struct {
	BackFill      bool
	BackFillOnly  bool
	Start         *big.Int
	End           *big.Int // set to 0 or a negative value to have no ending block
	HeaderFilter  HeaderFilter
	TxFilter      TxFilter
	ReceiptFilter ReceiptFilter
	StateFilter   StateFilter
	StorageFilter StorageFilter
}

// HeaderFilter contains filter settings for headers
type HeaderFilter struct {
	Off    bool
	Uncles bool
}

// TxFilter contains filter settings for txs
type TxFilter struct {
	Off bool
	Src []string
	Dst []string
}

// ReceiptFilter contains filter settings for receipts
type ReceiptFilter struct {
	Off       bool
	MatchTxs  bool // turn on to retrieve receipts that pair with retrieved transactions
	Contracts []string
	Topic0s   []string
}

// StateFilter contains filter settings for state
type StateFilter struct {
	Off               bool
	Addresses         []string // is converted to state key by taking its keccak256 hash
	IntermediateNodes bool
}

// StorageFilter contains filter settings for storage
type StorageFilter struct {
	Off               bool
	Addresses         []string
	StorageKeys       []string
	IntermediateNodes bool
}

// Init is used to initialize a EthSubscription struct with env variables
func NewEthSubscriptionConfig() *EthSubscription {
	sc := new(EthSubscription)
	// Below default to false, which means we do not backfill by default
	sc.BackFill = viper.GetBool("superNode.ethSubscription.historicalData")
	sc.BackFillOnly = viper.GetBool("superNode.ethSubscription.historicalDataOnly")
	// Below default to 0
	// 0 start means we start at the beginning and 0 end means we continue indefinitely
	sc.Start = big.NewInt(viper.GetInt64("superNode.ethSubscription.startingBlock"))
	sc.End = big.NewInt(viper.GetInt64("superNode.ethSubscription.endingBlock"))
	// Below default to false, which means we get all headers and no uncles by default
	sc.HeaderFilter = HeaderFilter{
		Off:    viper.GetBool("superNode.ethSubscription.off"),
		Uncles: viper.GetBool("superNode.ethSubscription.uncles"),
	}
	// Below defaults to false and two slices of length 0
	// Which means we get all transactions by default
	sc.TxFilter = TxFilter{
		Off: viper.GetBool("superNode.ethSubscription.trxFilter.off"),
		Src: viper.GetStringSlice("superNode.ethSubscription.trxFilter.src"),
		Dst: viper.GetStringSlice("superNode.ethSubscription.trxFilter.dst"),
	}
	// Below defaults to false and one slice of length 0
	// Which means we get all receipts by default
	sc.ReceiptFilter = ReceiptFilter{
		Off:       viper.GetBool("superNode.ethSubscription.receiptFilter.off"),
		Contracts: viper.GetStringSlice("superNode.ethSubscription.receiptFilter.contracts"),
		Topic0s:   viper.GetStringSlice("superNode.ethSubscription.receiptFilter.topic0s"),
	}
	// Below defaults to two false, and a slice of length 0
	// Which means we get all state leafs by default, but no intermediate nodes
	sc.StateFilter = StateFilter{
		Off:               viper.GetBool("superNode.ethSubscription.stateFilter.off"),
		IntermediateNodes: viper.GetBool("superNode.ethSubscription.stateFilter.intermediateNodes"),
		Addresses:         viper.GetStringSlice("superNode.ethSubscription.stateFilter.addresses"),
	}
	// Below defaults to two false, and two slices of length 0
	// Which means we get all storage leafs by default, but no intermediate nodes
	sc.StorageFilter = StorageFilter{
		Off:               viper.GetBool("superNode.ethSubscription.storageFilter.off"),
		IntermediateNodes: viper.GetBool("superNode.ethSubscription.storageFilter.intermediateNodes"),
		Addresses:         viper.GetStringSlice("superNode.ethSubscription.storageFilter.addresses"),
		StorageKeys:       viper.GetStringSlice("superNode.ethSubscription.storageFilter.storageKeys"),
	}
	return sc
}

// StartingBlock satisfies the SubscriptionSettings() interface
func (sc *EthSubscription) StartingBlock() *big.Int {
	return sc.Start
}

// EndingBlock satisfies the SubscriptionSettings() interface
func (sc *EthSubscription) EndingBlock() *big.Int {
	return sc.End
}

// HistoricalData satisfies the SubscriptionSettings() interface
func (sc *EthSubscription) HistoricalData() bool {
	return sc.BackFill
}

// HistoricalDataOnly satisfies the SubscriptionSettings() interface
func (sc *EthSubscription) HistoricalDataOnly() bool {
	return sc.BackFillOnly
}

// ChainType satisfies the SubscriptionSettings() interface
func (sc *EthSubscription) ChainType() ChainType {
	return Ethereum
}
