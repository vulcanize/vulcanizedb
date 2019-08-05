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
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// ISeedNodeStreamer is the interface for streaming data from a vulcanizeDB seed node
type ISeedNodeStreamer interface {
	Stream(payloadChan chan SeedNodePayload, streamFilters config.Subscription) (*rpc.ClientSubscription, error)
}

// SeedNodeStreamer is the underlying struct for the ISeedNodeStreamer interface
type SeedNodeStreamer struct {
	Client core.RpcClient
}

// NewSeedNodeStreamer creates a pointer to a new SeedNodeStreamer which satisfies the ISeedNodeStreamer interface
func NewSeedNodeStreamer(client core.RpcClient) *SeedNodeStreamer {
	return &SeedNodeStreamer{
		Client: client,
	}
}

// Stream is the main loop for subscribing to data from a vulcanizedb seed node
func (sds *SeedNodeStreamer) Stream(payloadChan chan SeedNodePayload, streamFilters config.Subscription) (*rpc.ClientSubscription, error) {
	return sds.Client.Subscribe("vulcanizedb", payloadChan, "stream", streamFilters)
}

// Payload holds the data returned from the seed node to the requesting client
type SeedNodePayload struct {
	BlockNumber     *big.Int                               `json:"blockNumber"`
	HeadersRlp      [][]byte                               `json:"headersRlp"`
	UnclesRlp       [][]byte                               `json:"unclesRlp"`
	TransactionsRlp [][]byte                               `json:"transactionsRlp"`
	ReceiptsRlp     [][]byte                               `json:"receiptsRlp"`
	StateNodesRlp   map[common.Hash][]byte                 `json:"stateNodesRlp"`
	StorageNodesRlp map[common.Hash]map[common.Hash][]byte `json:"storageNodesRlp"`
	ErrMsg          string                                 `json:"errMsg"`

	encoded []byte
	err     error
}

func (sd *SeedNodePayload) ensureEncoded() {
	if sd.encoded == nil && sd.err == nil {
		sd.encoded, sd.err = json.Marshal(sd)
	}
}

// Length to implement Encoder interface for StateDiff
func (sd *SeedNodePayload) Length() int {
	sd.ensureEncoded()
	return len(sd.encoded)
}

// Encode to implement Encoder interface for StateDiff
func (sd *SeedNodePayload) Encode() ([]byte, error) {
	sd.ensureEncoded()
	return sd.encoded, sd.err
}
