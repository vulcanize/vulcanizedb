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

package storage

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

// BackFiller is the backfilling interface
type BackFiller interface {
	BackFill(startingBlock, endingBlock uint64) ([]utils.StorageDiff, error)
}

// backFiller is the backfilling struct
type backFiller struct {
	fetcher fetcher.StateDiffFetcher
}

// NewStorageBackFiller returns a BackFiller
func NewStorageBackFiller(fetcher fetcher.StateDiffFetcher) BackFiller {
	return &backFiller{
		fetcher: fetcher,
	}
}

// BackFill uses the provided config to fetch and return the state diff at the specified blocknumber
// StateDiffAt(ctx context.Context, blockNumber uint64) (*Payload, error)
func (bf *backFiller) BackFill(startingBlock, endingBlock uint64) ([]utils.StorageDiff, error) {
	results := make([]utils.StorageDiff, 0)
	if endingBlock < startingBlock {
		return nil, errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	blockHeights := make([]uint64, 0, endingBlock-startingBlock+1)
	for i := startingBlock; i <= endingBlock; i++ {
		blockHeights = append(blockHeights, i)
	}
	payloads, err := bf.fetcher.FetchStateDiffsAt(blockHeights)
	if err != nil {
		return nil, err
	}
	for _, payload := range payloads {
		stateDiff := new(statediff.StateDiff)
		stateDiffDecodeErr := rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
		if stateDiffDecodeErr != nil {
			return nil, stateDiffDecodeErr
		}
		accounts := utils.GetAccountsFromDiff(*stateDiff)
		for _, account := range accounts {
			logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account", len(account.Storage)))
			for _, storage := range account.Storage {
				diff, formatErr := utils.FromGethStateDiff(account, stateDiff, storage)
				if formatErr != nil {
					return nil, formatErr
				}
				logrus.Trace("adding storage diff to results",
					"keccak of address: ", diff.HashedAddress.Hex(),
					"block height: ", diff.BlockHeight,
					"storage key: ", diff.StorageKey.Hex(),
					"storage value: ", diff.StorageValue.Hex())
				results = append(results, diff)
			}
		}
	}
	return results, nil
}
