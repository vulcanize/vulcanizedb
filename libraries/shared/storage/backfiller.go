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
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

const (
	DefaultMaxBatchSize   uint64 = 5000
	defaultMaxBatchNumber int64  = 100
)

// BackFiller is the backfilling interface
type BackFiller interface {
	BackFill(endingBlock uint64, backFill chan utils.StorageDiff, errChan chan error, done chan bool) error
}

// backFiller is the backfilling struct
type backFiller struct {
	fetcher       fetcher.StateDiffFetcher
	batchSize     uint64
	startingBlock uint64
}

// NewStorageBackFiller returns a BackFiller
func NewStorageBackFiller(fetcher fetcher.StateDiffFetcher, startingBlock, batchSize uint64) BackFiller {
	if batchSize == 0 {
		batchSize = DefaultMaxBatchSize
	}
	return &backFiller{
		fetcher:       fetcher,
		batchSize:     batchSize,
		startingBlock: startingBlock,
	}
}

// BackFill fetches, processes, and returns utils.StorageDiffs over a range of blocks
// It splits a large range up into smaller chunks, batch fetching and processing those chunks concurrently
func (bf *backFiller) BackFill(endingBlock uint64, backFill chan utils.StorageDiff, errChan chan error, done chan bool) error {
	if endingBlock < bf.startingBlock {
		return errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	// break the range up into bins of smaller ranges
	length := endingBlock - bf.startingBlock + 1
	numberOfBins := length / bf.batchSize
	remainder := length % bf.batchSize
	if remainder != 0 {
		numberOfBins++
	}
	blockRangeBins := make([][]uint64, numberOfBins)
	for i := range blockRangeBins {
		nextBinStart := bf.startingBlock + uint64(bf.batchSize)
		if nextBinStart > endingBlock {
			nextBinStart = endingBlock + 1
		}
		blockRange := make([]uint64, 0, nextBinStart-bf.startingBlock+1)
		for j := bf.startingBlock; j < nextBinStart; j++ {
			blockRange = append(blockRange, j)
		}
		bf.startingBlock = nextBinStart
		blockRangeBins[i] = blockRange
	}

	// int64 for atomic incrementing and decrementing to track the number of active processing goroutines we have
	var activeCount int64
	// channel for processing goroutines to signal when they are done
	processingDone := make(chan bool)

	// for each block range bin spin up a goroutine to batch fetch and process state diffs for that range
	go func() {
		for _, blockHeights := range blockRangeBins {
			// if we have reached our limit of active goroutines
			// wait for one to finish before starting the next
			if atomic.AddInt64(&activeCount, 1) > defaultMaxBatchNumber {
				// this blocks until a process signals it has finished
				// immediately forwards the signal to the normal listener so that it keeps the correct count
				processingDone <- <-processingDone
			}
			go func(blockHeights []uint64) {
				payloads, fetchErr := bf.fetcher.FetchStateDiffsAt(blockHeights)
				if fetchErr != nil {
					errChan <- fetchErr
				}
				for _, payload := range payloads {
					stateDiff := new(statediff.StateDiff)
					stateDiffDecodeErr := rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
					if stateDiffDecodeErr != nil {
						errChan <- stateDiffDecodeErr
						continue
					}
					accounts := utils.GetAccountsFromDiff(*stateDiff)
					for _, account := range accounts {
						logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account", len(account.Storage)))
						for _, storage := range account.Storage {
							diff, formatErr := utils.FromGethStateDiff(account, stateDiff, storage)
							if formatErr != nil {
								errChan <- formatErr
								continue
							}
							logrus.Trace("adding storage diff to results",
								"keccak of address: ", diff.HashedAddress.Hex(),
								"block height: ", diff.BlockHeight,
								"storage key: ", diff.StorageKey.Hex(),
								"storage value: ", diff.StorageValue.Hex())
							backFill <- diff
						}
					}
				}
				// when this goroutine is done, send out a signal
				processingDone <- true
			}(blockHeights)
		}
	}()

	// goroutine that listens on the processingDone chan
	// keeps track of the number of processing goroutines that have finished
	// when they have all finished, sends the final signal out
	go func() {
		goroutinesFinished := 0
		for {
			select {
			case <-processingDone:
				atomic.AddInt64(&activeCount, -1)
				goroutinesFinished++
				if goroutinesFinished == int(numberOfBins) {
					done <- true
					return
				}
			}
		}
	}()

	return nil
}
