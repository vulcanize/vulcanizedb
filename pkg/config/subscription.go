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

import "math/big"

// Subscription config is used by a subscribing transformer to specifiy which data to receive from the super node
type Subscription struct {
	BackFill      bool
	BackFillOnly  bool
	StartingBlock *big.Int
	EndingBlock   *big.Int // set to 0 or a negative value to have no ending block
	HeaderFilter  HeaderFilter
	TrxFilter     TrxFilter
	ReceiptFilter ReceiptFilter
	StateFilter   StateFilter
	StorageFilter StorageFilter
}

type HeaderFilter struct {
	Off    bool
	Uncles bool
}

type TrxFilter struct {
	Off bool
	Src []string
	Dst []string
}

type ReceiptFilter struct {
	Off       bool
	Contracts []string
	Topic0s   []string
}

type StateFilter struct {
	Off               bool
	Addresses         []string // is converted to state key by taking its keccak256 hash
	IntermediateNodes bool
}

type StorageFilter struct {
	Off               bool
	Addresses         []string
	StorageKeys       []string
	IntermediateNodes bool
}
