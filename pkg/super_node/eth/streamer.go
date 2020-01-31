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

package eth

import (
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
)

const (
	PayloadChanBufferSize = 20000 // the max eth sub buffer size
)

// PayloadStreamer satisfies the PayloadStreamer interface for ethereum
type PayloadStreamer struct {
	Client core.RPCClient
}

// NewPayloadStreamer creates a pointer to a new PayloadStreamer which satisfies the PayloadStreamer interface for ethereum
func NewPayloadStreamer(client core.RPCClient) *PayloadStreamer {
	return &PayloadStreamer{
		Client: client,
	}
}

// Stream is the main loop for subscribing to data from the Geth state diff process
// Satisfies the shared.PayloadStreamer interface
func (ps *PayloadStreamer) Stream(payloadChan chan shared.RawChainData) (shared.ClientSubscription, error) {
	logrus.Info("streaming diffs from geth")
	return ps.Client.Subscribe("statediff", payloadChan, "stream")
}
