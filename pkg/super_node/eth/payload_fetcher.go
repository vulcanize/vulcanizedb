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
	"fmt"

	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
)

// BatchClient is an interface to a batch-fetching geth rpc client; created to allow mock insertion
type BatchClient interface {
	BatchCall(batch []client.BatchElem) error
}

// PayloadFetcher satisfies the PayloadFetcher interface for ethereum
type PayloadFetcher struct {
	// PayloadFetcher is thread-safe as long as the underlying client is thread-safe, since it has/modifies no other state
	// http.Client is thread-safe
	client BatchClient
}

const method = "statediff_stateDiffAt"

// NewStateDiffFetcher returns a PayloadFetcher
func NewPayloadFetcher(bc BatchClient) *PayloadFetcher {
	return &PayloadFetcher{
		client: bc,
	}
}

// FetchAt fetches the statediff payloads at the given block heights
// Calls StateDiffAt(ctx context.Context, blockNumber uint64) (*Payload, error)
func (fetcher *PayloadFetcher) FetchAt(blockHeights []uint64) ([]interface{}, error) {
	batch := make([]client.BatchElem, 0)
	for _, height := range blockHeights {
		batch = append(batch, client.BatchElem{
			Method: method,
			Args:   []interface{}{height},
			Result: new(statediff.Payload),
		})
	}
	batchErr := fetcher.client.BatchCall(batch)
	if batchErr != nil {
		return nil, fmt.Errorf("PayloadFetcher err: %s", batchErr.Error())
	}
	results := make([]interface{}, 0, len(blockHeights))
	for _, batchElem := range batch {
		if batchElem.Error != nil {
			return nil, fmt.Errorf("PayloadFetcher err: %s", batchElem.Error.Error())
		}
		payload, ok := batchElem.Result.(*statediff.Payload)
		if ok {
			results = append(results, *payload)
		}
	}
	return results, nil
}
