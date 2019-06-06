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

// Streamer is used by watchers to stream eth data from a vulcanizedb seed node
package streamer

import (
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// IStreamer is the interface for streaming data from a vulcanizeDB seed node
type IStreamer interface {
	Stream(payloadChan chan ipfs.ResponsePayload, streamFilters ipfs.StreamFilters) (*rpc.ClientSubscription, error)
}

// Streamer is the underlying struct for the IStreamer interface
type Streamer struct {
	Client core.RpcClient
}

// NewSeedStreamer creates a pointer to a new Streamer which satisfies the IStreamer interface
func NewSeedStreamer(client core.RpcClient) *Streamer {
	return &Streamer{
		Client: client,
	}
}

// Stream is the main loop for subscribing to data from a vulcanizedb seed node
func (sds *Streamer) Stream(payloadChan chan ipfs.ResponsePayload, streamFilters ipfs.StreamFilters) (*rpc.ClientSubscription, error) {
	return sds.Client.Subscribe("vulcanizedb", payloadChan, "stream", streamFilters)
}
