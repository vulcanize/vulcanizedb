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

package super_node

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/config"
)

// Subscription holds the information for an individual client subscription to the super node
type Subscription struct {
	ID          rpc.ID
	PayloadChan chan<- Payload
	QuitChan    chan<- bool
}

// Payload is the struct for a super node stream payload
// It carries data of a type specific to the chain being supported/queried and an error message
type Payload struct {
	Data interface{} `json:"data"` // e.g. for Ethereum eth.StreamPayload
	Err  string      `json:"err"`
}

// SubscriptionSettings is the interface every subscription filter type needs to satisfy, no matter the chain
// Further specifics of the underlying filter type depend on the internal needs of the types
// which satisfy the ResponseFilterer and CIDRetriever interfaces for a specific chain
// The underlying type needs to be rlp serializable
type SubscriptionSettings interface {
	StartingBlock() *big.Int
	EndingBlock() *big.Int
	ChainType() config.ChainType
	HistoricalData() bool
	HistoricalDataOnly() bool
}
