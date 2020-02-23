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

// Streamer is used by watchers to stream eth data from a vulcanizedb super node
package streamer

import (
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// ISuperNodeStreamer is the interface for streaming SuperNodePayloads from a vulcanizeDB super node
type ISuperNodeStreamer interface {
	Stream(payloadChan chan super_node.SubscriptionPayload, params shared.SubscriptionSettings) (*rpc.ClientSubscription, error)
}

// SuperNodeStreamer is the underlying struct for the ISuperNodeStreamer interface
type SuperNodeStreamer struct {
	Client core.RPCClient
}

// NewSuperNodeStreamer creates a pointer to a new SuperNodeStreamer which satisfies the ISuperNodeStreamer interface
func NewSuperNodeStreamer(client core.RPCClient) *SuperNodeStreamer {
	return &SuperNodeStreamer{
		Client: client,
	}
}

// Stream is the main loop for subscribing to data from a vulcanizedb super node
func (sds *SuperNodeStreamer) Stream(payloadChan chan super_node.SubscriptionPayload, params shared.SubscriptionSettings) (*rpc.ClientSubscription, error) {
	return sds.Client.Subscribe("vdb", payloadChan, "stream", params)
}
