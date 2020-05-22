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
	"database/sql"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	storage2 "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

var ResultsLimit = 500

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
	AddTransformers(initializers []storage2.TransformerInitializer)
	Execute() error
}

type StorageWatcher struct {
	db                        *postgres.DB
	HeaderRepository          datastore.HeaderRepository
	KeccakAddressTransformers map[common.Hash]storage2.ITransformer // keccak hash of an address => transformer
	StorageDiffRepository     storage.DiffRepository
	DiffBlocksFromHeadOfChain int64 // the number of blocks from the head of the chain where diffs should be processed
}

func NewStorageWatcher(db *postgres.DB, backFromHeadOfChain int64) StorageWatcher {
	headerRepository := repositories.NewHeaderRepository(db)
	storageDiffRepository := storage.NewDiffRepository(db)
	transformers := make(map[common.Hash]storage2.ITransformer)
	return StorageWatcher{
		db:                        db,
		HeaderRepository:          headerRepository,
		KeccakAddressTransformers: transformers,
		StorageDiffRepository:     storageDiffRepository,
		DiffBlocksFromHeadOfChain: backFromHeadOfChain,
	}
}

func (watcher StorageWatcher) AddTransformers(initializers []storage2.TransformerInitializer) {
	for _, initializer := range initializers {
		storageTransformer := initializer(watcher.db)
		watcher.KeccakAddressTransformers[storageTransformer.KeccakContractAddress()] = storageTransformer
	}
}

func (watcher StorageWatcher) Execute() error {
	for {
		err := watcher.transformDiffs()
		if err != nil {
			logrus.Errorf("error transforming diffs: %s", err.Error())
			return err
		}
	}
}

func (watcher StorageWatcher) getMinDiffID() (int, error) {
	var minID = 0
	if watcher.DiffBlocksFromHeadOfChain != -1 {
		mostRecentHeaderBlockNumber, getHeaderErr := watcher.HeaderRepository.GetMostRecentHeaderBlockNumber()
		if getHeaderErr != nil {
			return 0, getHeaderErr
		}
		blockNumber := mostRecentHeaderBlockNumber - watcher.DiffBlocksFromHeadOfChain
		diffID, getDiffErr := watcher.StorageDiffRepository.GetFirstDiffIDForBlockHeight(blockNumber)
		if getDiffErr != nil {
			return 0, getDiffErr
		}

		// We are subtracting an offset from the diffID because it will be passed to GetNewDiffs which returns diffs with ids
		// greater than id passed in (minID), and we want to make sure that this diffID here is included in that collection
		diffOffset := int64(1)
		minID = int(diffID - diffOffset)
	}

	return minID, nil
}

func (watcher StorageWatcher) transformDiffs() error {
	minID, err := watcher.getMinDiffID()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	for {
		diffs, extractErr := watcher.StorageDiffRepository.GetNewDiffs(minID, ResultsLimit)
		if extractErr != nil {
			return fmt.Errorf("error getting unchecked diffs: %s", extractErr.Error())
		}
		for _, diff := range diffs {
			transformErr := watcher.transformDiff(diff)
			if transformErr != nil {
				if transformErr == sql.ErrNoRows || reflect.TypeOf(transformErr) == reflect.TypeOf(types.ErrKeyNotFound{}) {
					logrus.Tracef("error transforming diff: %s", transformErr.Error())
				} else {
					logrus.Infof("error transforming diff: %s", transformErr.Error())
				}
			}
		}
		lenDiffs := len(diffs)
		if lenDiffs > 0 {
			minID = int(diffs[lenDiffs-1].ID)
		}
		if lenDiffs < ResultsLimit {
			return nil
		}
	}
}

func (watcher StorageWatcher) transformDiff(diff types.PersistedDiff) error {
	t, watching := watcher.getTransformer(diff)
	if !watching {
		markCheckedErr := watcher.StorageDiffRepository.MarkChecked(diff.ID)
		if markCheckedErr != nil {
			return fmt.Errorf("error marking diff checked: %s", markCheckedErr.Error())
		}
		return nil
	}

	headerID, headerErr := watcher.getHeaderID(diff)
	if headerErr != nil {
		if headerErr == sql.ErrNoRows {
			return headerErr
		} else {
			return fmt.Errorf("error getting header for diff: %s", headerErr.Error())
		}
	}
	diff.HeaderID = headerID

	executeErr := t.Execute(diff)
	if executeErr != nil {
		if reflect.TypeOf(executeErr) == reflect.TypeOf(types.ErrKeyNotFound{}) {
			return executeErr
		} else {
			return fmt.Errorf("error executing storage transformer: %s", executeErr.Error())
		}
	}

	markCheckedErr := watcher.StorageDiffRepository.MarkChecked(diff.ID)
	if markCheckedErr != nil {
		return fmt.Errorf("error marking diff checked: %s", markCheckedErr.Error())
	}

	return nil
}

func (watcher StorageWatcher) getTransformer(diff types.PersistedDiff) (storage2.ITransformer, bool) {
	storageTransformer, ok := watcher.KeccakAddressTransformers[diff.HashedAddress]
	return storageTransformer, ok
}

func (watcher StorageWatcher) getHeaderID(diff types.PersistedDiff) (int64, error) {
	header, getHeaderErr := watcher.HeaderRepository.GetHeaderByBlockNumber(int64(diff.BlockHeight))
	if getHeaderErr != nil {
		return 0, getHeaderErr
	}
	if diff.BlockHash != common.HexToHash(header.Hash) {
		return 0, NewErrHeaderMismatch(header.Hash, diff.BlockHash.Hex())
	}
	return header.Id, nil
}
