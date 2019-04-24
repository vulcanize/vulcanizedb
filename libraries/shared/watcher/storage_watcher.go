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
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type StorageWatcher struct {
	db             *postgres.DB
	StorageFetcher fetcher.IStorageFetcher
	Queue          storage.IStorageQueue
	Transformers   map[common.Address]transformer.StorageTransformer
}

func NewStorageWatcher(fetcher fetcher.IStorageFetcher, db *postgres.DB) StorageWatcher {
	transformers := make(map[common.Address]transformer.StorageTransformer)
	queue := storage.NewStorageQueue(db)
	return StorageWatcher{
		db:             db,
		StorageFetcher: fetcher,
		Queue:          queue,
		Transformers:   transformers,
	}
}

func (watcher StorageWatcher) AddTransformers(initializers []transformer.StorageTransformerInitializer) {
	for _, initializer := range initializers {
		storageTransformer := initializer(watcher.db)
		watcher.Transformers[storageTransformer.ContractAddress()] = storageTransformer
	}
}

func (watcher StorageWatcher) Execute() error {
	rows := make(chan utils.StorageDiffRow)
	errs := make(chan error)
	go watcher.StorageFetcher.FetchStorageDiffs(rows, errs)
	for {
		select {
		case row := <-rows:
			watcher.processRow(row)
		case err := <-errs:
			return err
		}
	}
}

func (watcher StorageWatcher) processRow(row utils.StorageDiffRow) {
	storageTransformer, ok := watcher.Transformers[row.Contract]
	if !ok {
		// ignore rows from unwatched contracts
		return
	}
	executeErr := storageTransformer.Execute(row)
	if executeErr != nil {
		if isKeyNotFound(executeErr) {
			queueErr := watcher.Queue.Add(row)
			if queueErr != nil {
				logrus.Warn(queueErr.Error())
			}
		} else {
			logrus.Warn(executeErr.Error())
		}
	}
}

func isKeyNotFound(executeErr error) bool {
	return reflect.TypeOf(executeErr) == reflect.TypeOf(utils.ErrStorageKeyNotFound{})
}
