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

package shared

import (
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// Repository is the interface for the Postgres database
type Repository interface {
	LoadTriggers() error
	QueueData(payload super_node.SubscriptionPayload) error
	GetQueueData(height int64, hash string) (super_node.SubscriptionPayload, error)
	ReadyData(payload super_node.SubscriptionPayload) error
}

// SuperNodeStreamer is the interface for streaming data from a vulcanizeDB super node
type SuperNodeStreamer interface {
	Stream(payloadChan chan super_node.SubscriptionPayload, params shared.SubscriptionSettings) (*rpc.ClientSubscription, error)
}
