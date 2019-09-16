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
	"reflect"
	syn "sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	storageUtils "github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type StorageWatcherInterface interface {
	AddTransformers(initializers []transformer.StorageTransformerInitializer)
	Execute(rows chan utils.StorageDiffRow, errs chan error)
	WatchEthStorage(wg *syn.WaitGroup)
}

type StorageWatcher struct {
	db                   *postgres.DB
	StorageFetcher       fetcher.IStorageFetcher
	Queue                storage.IStorageQueue
	Transformers         map[common.Address]transformer.StorageTransformer
	PollingInterval      time.Duration
	QueueRecheckInterval time.Duration
}

func NewStorageWatcher(fetcher fetcher.IStorageFetcher, db *postgres.DB, pollingInterval time.Duration, queueRecheckInterval time.Duration) StorageWatcher {
	transformers := make(map[common.Address]transformer.StorageTransformer)
	queue := storage.NewStorageQueue(db)
	return StorageWatcher{
		db:                   db,
		StorageFetcher:       fetcher,
		Queue:                queue,
		Transformers:         transformers,
		PollingInterval:      pollingInterval,
		QueueRecheckInterval: queueRecheckInterval,
	}
}

func (storageWatcher StorageWatcher) AddTransformers(initializers []transformer.StorageTransformerInitializer) {
	for _, initializer := range initializers {
		storageTransformer := initializer(storageWatcher.db)
		storageWatcher.Transformers[storageTransformer.ContractAddress()] = storageTransformer
	}
}

func (storageWatcher StorageWatcher) Execute(rows chan utils.StorageDiffRow, errs chan error) {
	ticker := time.NewTicker(storageWatcher.QueueRecheckInterval)
	go storageWatcher.StorageFetcher.FetchStorageDiffs(rows, errs)
	for {
		select {
		case fetchErr := <-errs:
			logrus.Warn(fmt.Sprintf("error fetching storage diffs: %s", fetchErr))
		case row := <-rows:
			storageWatcher.processRow(row)
		case <-ticker.C:
			storageWatcher.processQueue()
		}
	}
}

func (storageWatcher StorageWatcher) WatchEthStorage(wg *syn.WaitGroup) {
	defer wg.Done()
	// Execute over the StorageTransformerInitializer set using the storage watcher
	ticker := time.NewTicker(storageWatcher.PollingInterval)
	defer ticker.Stop()
	for range ticker.C {
		errs := make(chan error)
		rows := make(chan storageUtils.StorageDiffRow)
		storageWatcher.Execute(rows, errs)
	}
}

func (storageWatcher StorageWatcher) processRow(row utils.StorageDiffRow) {
	storageTransformer, ok := storageWatcher.Transformers[row.Contract]
	if !ok {
		// ignore rows from unwatched contracts
		return
	}
	executeErr := storageTransformer.Execute(row)
	if executeErr != nil {
		logrus.Warn(fmt.Sprintf("error executing storage transformer: %s", executeErr))
		queueErr := storageWatcher.Queue.Add(row)
		if queueErr != nil {
			logrus.Warn(fmt.Sprintf("error queueing storage diff: %s", queueErr))
		}
	}
}

func (storageWatcher StorageWatcher) processQueue() {
	rows, fetchErr := storageWatcher.Queue.GetAll()
	if fetchErr != nil {
		logrus.Warn(fmt.Sprintf("error getting queued storage: %s", fetchErr))
	}
	for _, row := range rows {
		storageTransformer, ok := storageWatcher.Transformers[row.Contract]
		if !ok {
			// delete row from queue if address no longer watched
			storageWatcher.deleteRow(row.Id)
			continue
		}
		executeErr := storageTransformer.Execute(row)
		if executeErr == nil {
			storageWatcher.deleteRow(row.Id)
		}
	}
}

func (storageWatcher StorageWatcher) deleteRow(id int) {
	deleteErr := storageWatcher.Queue.Delete(id)
	if deleteErr != nil {
		logrus.Warn(fmt.Sprintf("error deleting persisted row from queue: %s", deleteErr))
	}
}

func isKeyNotFound(executeErr error) bool {
	return reflect.TypeOf(executeErr) == reflect.TypeOf(utils.ErrStorageKeyNotFound{})
}
