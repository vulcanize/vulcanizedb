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
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// BatchClient is an interface to a batch-fetching geth rpc client; created to allow mock insertion
type BatchClient interface {
	BatchCallContext(ctx context.Context, batch []rpc.BatchElem) error
}

// PayloadFetcher satisfies the PayloadFetcher interface for ethereum
type PayloadFetcher struct {
	// PayloadFetcher is thread-safe as long as the underlying client is thread-safe, since it has/modifies no other state
	// http.Client is thread-safe
	client  BatchClient
	timeout time.Duration
	params  statediff.Params
}

const method = "statediff_stateDiffAt"

// NewPayloadFetcher returns a PayloadFetcher
func NewPayloadFetcher(bc BatchClient, timeout time.Duration) *PayloadFetcher {
	return &PayloadFetcher{
		client:  bc,
		timeout: timeout,
		params: statediff.Params{
			IncludeReceipts:          true,
			IncludeTD:                true,
			IncludeBlock:             true,
			IntermediateStateNodes:   true,
			IntermediateStorageNodes: true,
		},
	}
}

// FetchAt fetches the statediff payloads at the given block heights
// Calls StateDiffAt(ctx context.Context, blockNumber uint64, params Params) (*Payload, error)
func (fetcher *PayloadFetcher) FetchAt(blockHeights []uint64) ([]shared.RawChainData, error) {
	batch := make([]rpc.BatchElem, 0)
	for _, height := range blockHeights {
		batch = append(batch, rpc.BatchElem{
			Method: method,
			Args:   []interface{}{height, fetcher.params},
			Result: new(statediff.Payload),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), fetcher.timeout)
	defer cancel()
	if err := fetcher.client.BatchCallContext(ctx, batch); err != nil {
		return nil, fmt.Errorf("ethereum PayloadFetcher batch err for block range %d-%d: %s", blockHeights[0], blockHeights[len(blockHeights)-1], err.Error())
	}
	results := make([]shared.RawChainData, 0, len(blockHeights))
	for _, batchElem := range batch {
		if batchElem.Error != nil {
			return nil, fmt.Errorf("ethereum PayloadFetcher err at blockheight %d: %s", batchElem.Args[0].(uint64), batchElem.Error.Error())
		}
		payload, ok := batchElem.Result.(*statediff.Payload)
		if ok {
			results = append(results, *payload)
		}
	}
	return results, nil
}
