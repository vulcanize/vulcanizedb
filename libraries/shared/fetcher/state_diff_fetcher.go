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

package fetcher

import (
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
)

// IStateDiffFetcher is the state diff fetching interface
type IStateDiffFetcher interface {
	FetchStateDiffsAt(blockHeights []uint64) ([]*statediff.Payload, error)
}

// BatchClient is an interface to a batch-fetching geth rpc client; created to allow mock insertion
type BatchClient interface {
	BatchCall(batch []client.BatchElem) error
}

// StateDiffFetcher is the state diff fetching struct
type StateDiffFetcher struct {
	client BatchClient
}

const method = "statediff_stateDiffAt"

// NewStateDiffFetcher returns a IStateDiffFetcher
func NewStateDiffFetcher(bc BatchClient) IStateDiffFetcher {
	return &StateDiffFetcher{
		client: bc,
	}
}

// FetchStateDiffsAt fetches the statediff payloads at the given block heights
// Calls StateDiffAt(ctx context.Context, blockNumber uint64) (*Payload, error)
func (sdf *StateDiffFetcher) FetchStateDiffsAt(blockHeights []uint64) ([]*statediff.Payload, error) {
	batch := make([]client.BatchElem, 0)
	for _, height := range blockHeights {
		batch = append(batch, client.BatchElem{
			Method: method,
			Args:   []interface{}{height},
			Result: new(statediff.Payload),
		})
	}
	batchErr := sdf.client.BatchCall(batch)
	if batchErr != nil {
		return nil, batchErr
	}
	results := make([]*statediff.Payload, 0, len(blockHeights))
	for _, batchElem := range batch {
		if batchElem.Error != nil {
			return nil, batchElem.Error
		}
		results = append(results, batchElem.Result.(*statediff.Payload))
	}
	return results, nil
}
