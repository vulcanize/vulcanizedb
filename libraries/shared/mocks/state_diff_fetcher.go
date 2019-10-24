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

	"github.com/ethereum/go-ethereum/statediff"
)

// StateDiffFetcher mock for tests
type StateDiffFetcher struct {
	PayloadsToReturn     map[uint64]*statediff.Payload
	FetchErr             error
	CalledAtBlockHeights [][]uint64
}

// SetPayloadsToReturn for tests
func (fetcher *StateDiffFetcher) SetPayloadsToReturn(payloads map[uint64]*statediff.Payload) {
	fetcher.PayloadsToReturn = payloads
}

// FetchStateDiffsAt mock method
func (fetcher *StateDiffFetcher) FetchStateDiffsAt(blockHeights []uint64) ([]*statediff.Payload, error) {
	fetcher.CalledAtBlockHeights = append(fetcher.CalledAtBlockHeights, blockHeights)
	if fetcher.PayloadsToReturn == nil {
		return nil, errors.New("MockStateDiffFetcher needs to be initialized with payloads to return")
	}
	results := make([]*statediff.Payload, 0, len(blockHeights))
	for _, height := range blockHeights {
		results = append(results, fetcher.PayloadsToReturn[height])
	}
	return results, nil
}
