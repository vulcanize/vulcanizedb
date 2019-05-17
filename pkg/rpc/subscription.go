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

package rpc

// Subscription
type Subscription struct {
}

// Params are set by the client to tell the server how to filter that is fed into their subscription
type Params struct {
	HeaderFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64 // set to 0 or a negative value to have no ending block
		Uncles        bool
	}
	TrxFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Src           string
		Dst           string
	}
	ReceiptFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Topic0s       []string
	}
	StateFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Address       string // is converted to state key by taking its keccak256 hash
		LeafsOnly     bool
	}
	StorageFilter struct {
		Off           bool
		StartingBlock int64
		EndingBlock   int64
		Address       string
		StorageKey    string
		LeafsOnly     bool
	}
}
