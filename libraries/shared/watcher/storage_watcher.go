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

package watcher

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

type IStorageWatcher interface {
	AddTransformers(initializers []transformer.StorageTransformerInitializer)
	Execute(queueRecheckInterval time.Duration, backFillOn bool)
	BackFill(startingBlock uint64, backFiller storage.BackFiller)
}

type StorageWatcher struct {
	db                        *postgres.DB
	StorageFetcher            fetcher.IStorageFetcher
	Queue                     storage.IStorageQueue
	StorageDiffRepository     datastore.StorageDiffRepository
	KeccakAddressTransformers map[common.Hash]transformer.StorageTransformer // keccak hash of an address => transformer
	DiffsChan                 chan utils.StorageDiffInput
	ErrsChan                  chan error
	BackFillDoneChan          chan bool
	StartingSyncBlockChan     chan uint64
}

func NewStorageWatcher(f fetcher.IStorageFetcher, db *postgres.DB) *StorageWatcher {
	queue := storage.NewStorageQueue(db)
	storageDiffRepository := repositories.NewStorageDiffRepository(db)
	transformers := make(map[common.Hash]transformer.StorageTransformer)
	return &StorageWatcher{
		db:                        db,
		StorageFetcher:            f,
		DiffsChan:                 make(chan utils.StorageDiffInput, fetcher.PayloadChanBufferSize),
		ErrsChan:                  make(chan error),
		StartingSyncBlockChan:     make(chan uint64),
		BackFillDoneChan:          make(chan bool),
		Queue:                     queue,
		StorageDiffRepository:     storageDiffRepository,
		KeccakAddressTransformers: transformers,
	}
}

func (storageWatcher *StorageWatcher) AddTransformers(initializers []transformer.StorageTransformerInitializer) {
	for _, initializer := range initializers {
		storageTransformer := initializer(storageWatcher.db)
		storageWatcher.KeccakAddressTransformers[storageTransformer.KeccakContractAddress()] = storageTransformer
	}
}

// BackFill uses a backFiller to backfill missing storage diffs for the storageWatcher
func (storageWatcher *StorageWatcher) BackFill(startingBlock uint64, backFiller storage.BackFiller) {
	// this blocks until the Execute process sends us the first block number it sees
	endBackFillBlock := <-storageWatcher.StartingSyncBlockChan
	backFillInitErr := backFiller.BackFill(startingBlock, endBackFillBlock,
		storageWatcher.DiffsChan, storageWatcher.ErrsChan, storageWatcher.BackFillDoneChan)
	if backFillInitErr != nil {
		logrus.Warn(backFillInitErr)
	}
}

// Execute runs the StorageWatcher processes
func (storageWatcher *StorageWatcher) Execute(queueRecheckInterval time.Duration, backFillOn bool) {
	ticker := time.NewTicker(queueRecheckInterval)
	go storageWatcher.StorageFetcher.FetchStorageDiffs(storageWatcher.DiffsChan, storageWatcher.ErrsChan)
	start := true
	for {
		select {
		case err := <-storageWatcher.ErrsChan:
			logrus.Warn(fmt.Sprintf("error fetching storage diffs: %s", err.Error()))
		case diff := <-storageWatcher.DiffsChan:
			if start && backFillOn {
				storageWatcher.StartingSyncBlockChan <- uint64(diff.BlockHeight - 1)
				start = false
			}
			storageWatcher.processRow(diff)
		case <-ticker.C:
			storageWatcher.processQueue()
		case <-storageWatcher.BackFillDoneChan:
			logrus.Info("storage watcher backfill process has finished")
		}
	}
}

func (storageWatcher *StorageWatcher) getTransformer(diff utils.PersistedStorageDiff) (transformer.StorageTransformer, bool) {
	storageTransformer, ok := storageWatcher.KeccakAddressTransformers[diff.HashedAddress]
	return storageTransformer, ok
}

func (storageWatcher StorageWatcher) processRow(diffInput utils.StorageDiffInput) {
	diffID, err := storageWatcher.StorageDiffRepository.CreateStorageDiff(diffInput)
	if err != nil {
		if err == repositories.ErrDuplicateDiff {
			logrus.Warn("ignoring duplicate diff")
			return
		}
		logrus.Warnf("failed to persist storage diff: %s", err.Error())
		// TODO: bail? Should we continue attempting to transform a diff we didn't persist
	}
	persistedDiff := utils.ToPersistedDiff(diffInput, diffID)
	storageTransformer, ok := storageWatcher.getTransformer(persistedDiff)
	if !ok {
		logrus.Debug("ignoring diff from unwatched contract")
		return
	}
	executeErr := storageTransformer.Execute(persistedDiff)
	if executeErr != nil {
		logrus.Warn(fmt.Sprintf("error executing storage transformer: %s", executeErr))
		queueErr := storageWatcher.Queue.Add(persistedDiff)
		if queueErr != nil {
			logrus.Warn(fmt.Sprintf("error queueing storage diff: %s", queueErr))
		}
		return
	}
	logrus.Debugf("Storage diff persisted at block height: %d", diffInput.BlockHeight)
}

func (storageWatcher StorageWatcher) processQueue() {
	diffs, fetchErr := storageWatcher.Queue.GetAll()
	if fetchErr != nil {
		logrus.Warn(fmt.Sprintf("error getting queued storage: %s", fetchErr))
	}
	for _, diff := range diffs {
		storageTransformer, ok := storageWatcher.getTransformer(diff)
		if !ok {
			// delete diff from queue if address no longer watched
			storageWatcher.deleteRow(diff.ID)
			continue
		}
		executeErr := storageTransformer.Execute(diff)
		if executeErr == nil {
			storageWatcher.deleteRow(diff.ID)
		}
	}
}

func (storageWatcher StorageWatcher) deleteRow(diffID int64) {
	deleteErr := storageWatcher.Queue.Delete(diffID)
	if deleteErr != nil {
		logrus.Warn(fmt.Sprintf("error deleting persisted diff from queue: %s", deleteErr))
	}
}
