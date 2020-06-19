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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	storage2 "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fs"
	"github.com/sirupsen/logrus"
)

var ResultsLimit = 500

var ErrHeaderMismatch = errors.New("header hash doesn't match between db and diff")

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
	StatusWriter              fs.StatusWriter
}

func NewStorageWatcher(db *postgres.DB, backFromHeadOfChain int64, statusWriter fs.StatusWriter) StorageWatcher {
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
	writeErr := watcher.StatusWriter.Write()
	if writeErr != nil {
		return fmt.Errorf("error confirming health check: %w", writeErr)
	}

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
			return 0, fmt.Errorf("error getting most recent header block number: %w", getHeaderErr)
		}
		blockNumber := mostRecentHeaderBlockNumber - watcher.DiffBlocksFromHeadOfChain
		diffID, getDiffErr := watcher.StorageDiffRepository.GetFirstDiffIDForBlockHeight(blockNumber)
		if getDiffErr != nil {
			return 0, fmt.Errorf("error getting first diff id for block height %d: %w", blockNumber, getDiffErr)
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error getting min diff ID: %w", err)
	}

	for {
		diffs, extractErr := watcher.StorageDiffRepository.GetNewDiffs(minID, ResultsLimit)
		if extractErr != nil {
			return fmt.Errorf("error getting unchecked diffs: %w", extractErr)
		}
		for _, diff := range diffs {
			transformErr := watcher.transformDiff(diff)
			if transformErr != nil {
				if errors.Is(transformErr, sql.ErrNoRows) || errors.Is(transformErr, types.ErrKeyNotFound) {
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
			return fmt.Errorf("error marking diff checked: %w", markCheckedErr)
		}
		return nil
	}

	headerID, headerErr := watcher.getHeaderID(diff)
	if headerErr != nil {
		return fmt.Errorf("error getting header for diff: %w", headerErr)
	}
	diff.HeaderID = headerID

	executeErr := t.Execute(diff)
	if executeErr != nil {
		return fmt.Errorf("error executing storage transformer: %w", executeErr)
	}

	markCheckedErr := watcher.StorageDiffRepository.MarkChecked(diff.ID)
	if markCheckedErr != nil {
		return fmt.Errorf("error marking diff checked: %w", markCheckedErr)
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
		return 0, fmt.Errorf("error getting header by block number %d: %w", diff.BlockHeight, getHeaderErr)
	}
	if diff.BlockHash != common.HexToHash(header.Hash) {
		msgToFormat := "diff ID %d, block %d, db hash %s, diff hash %s"
		details := fmt.Sprintf(msgToFormat, diff.ID, diff.BlockHeight, header.Hash, diff.BlockHash.Hex())
		return 0, fmt.Errorf("%w: %s", ErrHeaderMismatch, details)
	}
	return header.Id, nil
}
