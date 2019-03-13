// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
)

type StorageWatcher struct {
	db           *postgres.DB
	tailer       fs.Tailer
	Queue        storage.IStorageQueue
	Transformers map[common.Address]transformer.StorageTransformer
}

func NewStorageWatcher(tailer fs.Tailer, db *postgres.DB) StorageWatcher {
	transformers := make(map[common.Address]transformer.StorageTransformer)
	queue := storage.NewStorageQueue(db)
	return StorageWatcher{
		db:           db,
		tailer:       tailer,
		Queue:        queue,
		Transformers: transformers,
	}
}

func (watcher StorageWatcher) AddTransformers(initializers []transformer.StorageTransformerInitializer) {
	for _, initializer := range initializers {
		storageTransformer := initializer(watcher.db)
		watcher.Transformers[storageTransformer.ContractAddress()] = storageTransformer
	}
}

func (watcher StorageWatcher) Execute() error {
	t, tailErr := watcher.tailer.Tail()
	if tailErr != nil {
		return tailErr
	}
	for line := range t.Lines {
		row, parseErr := utils.FromStrings(strings.Split(line.Text, ","))
		if parseErr != nil {
			return parseErr
		}
		storageTransformer, ok := watcher.Transformers[row.Contract]
		if !ok {
			logrus.Warn(utils.ErrContractNotFound{Contract: row.Contract.Hex()}.Error())
			continue
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
			continue
		}
	}
	return nil
}

func isKeyNotFound(executeErr error) bool {
	return reflect.TypeOf(executeErr) == reflect.TypeOf(utils.ErrStorageKeyNotFound{})
}
