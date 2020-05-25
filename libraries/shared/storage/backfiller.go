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
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/utilities"
)

const (
	DefaultMaxBatchSize   uint64 = 100
	defaultMaxBatchNumber int64  = 10
)

// BackFiller is the backfilling interface
type BackFiller interface {
	BackFill(startingBlock, endingBlock uint64, backFill chan utils.StorageDiffInput, errChan chan error, done chan bool) error
}

// backFiller is the backfilling struct
type backFiller struct {
	fetcher       fetcher.StateDiffFetcher
	batchSize     uint64
	startingBlock uint64
}

// NewStorageBackFiller returns a BackFiller
func NewStorageBackFiller(fetcher fetcher.StateDiffFetcher, batchSize uint64) BackFiller {
	if batchSize == 0 {
		batchSize = DefaultMaxBatchSize
	}
	return &backFiller{
		fetcher:   fetcher,
		batchSize: batchSize,
	}
}

// BackFill fetches, processes, and returns utils.StorageDiffs over a range of blocks
// It splits a large range up into smaller chunks, batch fetching and processing those chunks concurrently
func (bf *backFiller) BackFill(startingBlock, endingBlock uint64, backFill chan utils.StorageDiffInput, errChan chan error, done chan bool) error {
	logrus.Infof("going to fill in gap from %d to %d", startingBlock, endingBlock)

	// break the range up into bins of smaller ranges
	blockRangeBins, err := utilities.GetBlockHeightBins(startingBlock, endingBlock, bf.batchSize)
	if err != nil {
		return err
	}
	// int64 for atomic incrementing and decrementing to track the number of active processing goroutines we have
	var activeCount int64
	// channel for processing goroutines to signal when they are done
	processingDone := make(chan [2]uint64)
	forwardDone := make(chan bool)

	// for each block range bin spin up a goroutine to batch fetch and process state diffs in that range
	go func() {
		for _, blockHeights := range blockRangeBins {
			// if we have reached our limit of active goroutines
			// wait for one to finish before starting the next
			if atomic.AddInt64(&activeCount, 1) > defaultMaxBatchNumber {
				// this blocks until a process signals it has finished
				<-forwardDone
			}
			go bf.backFillRange(blockHeights, backFill, errChan, processingDone)
		}
	}()

	// goroutine that listens on the processingDone chan
	// keeps track of the number of processing goroutines that have finished
	// when they have all finished, sends the final signal out
	go func() {
		goroutinesFinished := 0
		for {
			select {
			case doneWithHeights := <-processingDone:
				atomic.AddInt64(&activeCount, -1)
				select {
				// if we are waiting for a process to finish, signal that one has
				case forwardDone <- true:
				default:
				}
				logrus.Infof("finished fetching gap sub-bin from %d to %d", doneWithHeights[0], doneWithHeights[1])
				goroutinesFinished++
				if goroutinesFinished >= len(blockRangeBins) {
					done <- true
					return
				}
			}
		}
	}()

	return nil
}

func (bf *backFiller) backFillRange(blockHeights []uint64, diffChan chan utils.StorageDiffInput, errChan chan error, doneChan chan [2]uint64) {
	payloads, fetchErr := bf.fetcher.FetchStateDiffsAt(blockHeights)
	if fetchErr != nil {
		errChan <- fetchErr
	}
	for _, payload := range payloads {
		stateDiff := new(statediff.StateObject)
		stateDiffDecodeErr := rlp.DecodeBytes(payload.StateObjectRlp, stateDiff)
		if stateDiffDecodeErr != nil {
			errChan <- stateDiffDecodeErr
			continue
		}
		accounts := utils.GetAccountsFromDiff(*stateDiff)
		for _, account := range accounts {
			logrus.Trace(fmt.Sprintf("iterating through %d Storage values on account with key %s", len(account.StorageNodes), common.BytesToHash(account.LeafKey).Hex()))
			for _, storage := range account.StorageNodes {
				diff, formatErr := utils.FromGethStateDiff(account, stateDiff, storage)
				if formatErr != nil {
					logrus.Error("failed to format utils.StorageDiff from storage with key: ", common.BytesToHash(storage.LeafKey), "from account with key: ", common.BytesToHash(account.LeafKey))
					errChan <- formatErr
					continue
				}
				logrus.Trace("adding storage diff to results",
					"keccak of address: ", diff.HashedAddress.Hex(),
					"block height: ", diff.BlockHeight,
					"storage key: ", diff.StorageKey.Hex(),
					"storage value: ", diff.StorageValue.Hex())
				diffChan <- diff
			}
		}
	}
	// when this is done, send out a signal
	doneChan <- [2]uint64{blockHeights[0], blockHeights[len(blockHeights)-1]}
}
