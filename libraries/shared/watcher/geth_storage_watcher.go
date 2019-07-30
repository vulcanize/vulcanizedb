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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type GethStorageWatcher struct {
	StorageWatcher
}

func NewGethStorageWatcher(fetcher fetcher.IStorageFetcher, db *postgres.DB) GethStorageWatcher {
	queue := storage.NewStorageQueue(db)
	transformers := make(map[common.Address]transformer.StorageTransformer)
	keccakAddressTransformers := make(map[common.Address]transformer.StorageTransformer)
	storageWatcher := StorageWatcher{
		db:                        db,
		StorageFetcher:            fetcher,
		Queue:                     queue,
		Transformers:              transformers,
		KeccakAddressTransformers: keccakAddressTransformers,
	}
	storageWatcher.transformerGetter = storageWatcher.getTransformerForGethWatcher
	return GethStorageWatcher{StorageWatcher: storageWatcher}
}

func (storageWatcher StorageWatcher) getTransformerForGethWatcher(contractAddress common.Address) (transformer.StorageTransformer, bool) {
	storageTransformer, ok := storageWatcher.KeccakAddressTransformers[contractAddress]
	if ok {
		return storageTransformer, ok
	} else {
		for address, transformer := range storageWatcher.Transformers {
			keccakOfTransformerAddress := common.BytesToAddress(crypto.Keccak256(address[:]))
			if keccakOfTransformerAddress == contractAddress {
				storageWatcher.KeccakAddressTransformers[contractAddress] = transformer
				return transformer, true
			}
		}
	}

	return nil, false
}
