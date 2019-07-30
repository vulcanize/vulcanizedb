// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package watcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type CsvStorageWatcher struct {
	StorageWatcher
}

func NewCsvStorageWatcher(fetcher fetcher.IStorageFetcher, db *postgres.DB) CsvStorageWatcher {
	queue := storage.NewStorageQueue(db)
	transformers := make(map[common.Address]transformer.StorageTransformer)
	storageWatcher := StorageWatcher{
		db:             db,
		StorageFetcher: fetcher,
		Queue:          queue,
		Transformers:   transformers,
	}
	storageWatcher.transformerGetter = storageWatcher.getCsvTransformer
	return CsvStorageWatcher{StorageWatcher: storageWatcher}
}

func (storageWatcher StorageWatcher) getCsvTransformer(contractAddress common.Address) (transformer.StorageTransformer, bool) {
	storageTransformer, ok := storageWatcher.Transformers[contractAddress]
	return storageTransformer, ok
}
