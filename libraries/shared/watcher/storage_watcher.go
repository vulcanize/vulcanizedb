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
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

type ErrHeaderMismatch struct {
	dbHash   string
	diffHash string
}

func NewErrHeaderMismatch(DBHash, diffHash string) *ErrHeaderMismatch {
	return &ErrHeaderMismatch{dbHash: DBHash, diffHash: diffHash}
}

func (e ErrHeaderMismatch) Error() string {
	return fmt.Sprintf("db header hash (%s) doesn't match diff header hash (%s)", e.dbHash, e.diffHash)
}

type IStorageWatcher interface {
	AddTransformers(initializers []transformer.StorageTransformerInitializer)
	Execute(queueRecheckInterval time.Duration) error
}

type StorageWatcher struct {
	db                        *postgres.DB
	StorageFetcher            fetcher.IStorageFetcher
	Queue                     storage.IStorageQueue
	HeaderRepository          datastore.HeaderRepository
	StorageDiffRepository     datastore.StorageDiffRepository
	KeccakAddressTransformers map[common.Hash]transformer.StorageTransformer // keccak hash of an address => transformer
}

func NewStorageWatcher(fetcher fetcher.IStorageFetcher, db *postgres.DB) StorageWatcher {
	queue := storage.NewStorageQueue(db)
	headerRepository := repositories.NewHeaderRepository(db)
	storageDiffRepository := repositories.NewStorageDiffRepository(db)
	transformers := make(map[common.Hash]transformer.StorageTransformer)
	return StorageWatcher{
		db:                        db,
		StorageFetcher:            fetcher,
		Queue:                     queue,
		HeaderRepository:          headerRepository,
		StorageDiffRepository:     storageDiffRepository,
		KeccakAddressTransformers: transformers,
	}
}

func (storageWatcher StorageWatcher) AddTransformers(initializers []transformer.StorageTransformerInitializer) {
	for _, initializer := range initializers {
		storageTransformer := initializer(storageWatcher.db)
		storageWatcher.KeccakAddressTransformers[storageTransformer.KeccakContractAddress()] = storageTransformer
	}
}

func (storageWatcher StorageWatcher) Execute(queueRecheckInterval time.Duration) error {
	ticker := time.NewTicker(queueRecheckInterval)
	diffsChan := make(chan storage.RawStorageDiff)
	errsChan := make(chan error)

	defer close(diffsChan)
	defer close(errsChan)

	go storageWatcher.StorageFetcher.FetchStorageDiffs(diffsChan, errsChan)

	for {
		select {
		case fetchErr := <-errsChan:
			logrus.Warnf("error fetching storage diffs: %s", fetchErr.Error())
			return fetchErr
		case diff := <-diffsChan:
			storageWatcher.processRow(diff)
		case <-ticker.C:
			storageWatcher.processQueue()
		}
	}
}

func (storageWatcher StorageWatcher) getTransformer(diff storage.PersistedStorageDiff) (transformer.StorageTransformer, bool) {
	storageTransformer, ok := storageWatcher.KeccakAddressTransformers[diff.HashedAddress]
	return storageTransformer, ok
}

func (storageWatcher StorageWatcher) processRow(rawDiff storage.RawStorageDiff) {
	diffID, err := storageWatcher.StorageDiffRepository.CreateStorageDiff(rawDiff)
	if err != nil {
		if err == repositories.ErrDuplicateDiff {
			logrus.Info("ignoring duplicate diff")
			return
		}
		logrus.Warnf("failed to persist storage diff: %s", err.Error())
		// TODO: bail? Should we continue attempting to transform a diff we didn't persist
	}
	persistedDiff := storage.ToPersistedDiff(rawDiff, diffID)

	storageTransformer, isTransformerWatchingAddress := storageWatcher.getTransformer(persistedDiff)
	if !isTransformerWatchingAddress {
		logrus.Debug("ignoring diff from an unwatched contract")
		return
	}

	headerID, err := storageWatcher.getHeaderID(persistedDiff)
	if err != nil {
		logrus.Infof("error getting header for diff: %s", err.Error())
		storageWatcher.queueDiff(persistedDiff)
		return
	}
	persistedDiff.HeaderID = headerID

	executeErr := storageTransformer.Execute(persistedDiff)
	if executeErr != nil {
		logrus.Infof("error executing storage transformer: %s", executeErr.Error())
		storageWatcher.queueDiff(persistedDiff)
	}
}

func (storageWatcher StorageWatcher) processQueue() {
	diffs, fetchErr := storageWatcher.Queue.GetAll()
	if fetchErr != nil {
		logrus.Infof("error getting queued storage: %s", fetchErr.Error())
	}

	for _, diff := range diffs {
		storageTransformer, isTransformerWatchingAddress := storageWatcher.getTransformer(diff)
		if !isTransformerWatchingAddress {
			storageWatcher.deleteRow(diff.ID)
			continue
		}

		headerID, getHeaderErr := storageWatcher.getHeaderID(diff)
		if getHeaderErr != nil {
			logrus.Infof("error getting header for diff: %s", getHeaderErr.Error())
			continue
		}
		diff.HeaderID = headerID

		executeErr := storageTransformer.Execute(diff)
		if executeErr != nil {
			logrus.Infof("error executing storage transformer: %s", executeErr.Error())
			continue
		}

		storageWatcher.deleteRow(diff.ID)
	}
}

func (storageWatcher StorageWatcher) deleteRow(diffID int64) {
	deleteErr := storageWatcher.Queue.Delete(diffID)
	if deleteErr != nil {
		logrus.Infof("error deleting persisted diff from queue: %s", deleteErr.Error())
	}
}

func (storageWatcher StorageWatcher) queueDiff(diff storage.PersistedStorageDiff) {
	queueErr := storageWatcher.Queue.Add(diff)
	if queueErr != nil {
		logrus.Infof("error queueing storage diff: %s", queueErr.Error())
	}
}

func (storageWatcher StorageWatcher) getHeaderID(diff storage.PersistedStorageDiff) (int64, error) {
	header, getHeaderErr := storageWatcher.HeaderRepository.GetHeader(int64(diff.BlockHeight))
	if getHeaderErr != nil {
		return 0, getHeaderErr
	}
	if diff.BlockHash != common.HexToHash(header.Hash) {
		return 0, NewErrHeaderMismatch(header.Hash, diff.BlockHash.Hex())
	}
	return header.Id, nil
}
