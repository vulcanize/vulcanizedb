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

package streamer

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// IStateDiffStreamer is the interface for streaming a statediff subscription
type IStateDiffStreamer interface {
	Stream(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error)
}

// StateDiffStreamer is the underlying struct for the StateDiffStreamer interface
type StateDiffStreamer struct {
	Client core.RpcClient
}

// NewStateDiffStreamer creates a pointer to a new StateDiffStreamer which satisfies the IStateDiffStreamer interface
func NewStateDiffStreamer(client core.RpcClient) *StateDiffStreamer {
	return &StateDiffStreamer{
		Client: client,
	}
}

// Stream is the main loop for subscribing to data from the Geth state diff process
func (sds *StateDiffStreamer) Stream(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error) {
	return sds.Client.Subscribe("statediff", payloadChan, "stream")
}
