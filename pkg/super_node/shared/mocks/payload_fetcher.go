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

package mocks

import (
	"errors"
	"sync/atomic"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// PayloadFetcher mock for tests
type PayloadFetcher struct {
	PayloadsToReturn     map[uint64]shared.RawChainData
	FetchErrs            map[uint64]error
	CalledAtBlockHeights [][]uint64
	CalledTimes          int64
}

// FetchAt mock method
func (fetcher *PayloadFetcher) FetchAt(blockHeights []uint64) ([]shared.RawChainData, error) {
	if fetcher.PayloadsToReturn == nil {
		return nil, errors.New("mock StateDiffFetcher needs to be initialized with payloads to return")
	}
	atomic.AddInt64(&fetcher.CalledTimes, 1) // thread-safe increment
	fetcher.CalledAtBlockHeights = append(fetcher.CalledAtBlockHeights, blockHeights)
	results := make([]shared.RawChainData, 0, len(blockHeights))
	for _, height := range blockHeights {
		results = append(results, fetcher.PayloadsToReturn[height])
		err, ok := fetcher.FetchErrs[height]
		if ok && err != nil {
			return nil, err
		}
	}
	return results, nil
}
