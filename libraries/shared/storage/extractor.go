package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/makerdao/vulcanizedb/libraries/shared/storage/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

type DiffExtractor struct {
	StorageDiffRepository DiffRepository
	StorageFetcher        fetcher.IStorageFetcher
}

func NewDiffExtractor(fetcher fetcher.IStorageFetcher, db *postgres.DB) DiffExtractor {
	repo := NewDiffRepository(db)
	return DiffExtractor{
		StorageDiffRepository: repo,
		StorageFetcher:        fetcher,
	}
}

func (extractor DiffExtractor) ExtractDiffs() error {
	diffsChan := make(chan types.RawDiff)
	errsChan := make(chan error)

	defer close(diffsChan)
	defer close(errsChan)

	go extractor.StorageFetcher.FetchStorageDiffs(diffsChan, errsChan)

	for {
		select {
		case fetchErr := <-errsChan:
			logrus.Warnf("error fetching storage diffs: %s", fetchErr.Error())
			return fmt.Errorf("error fetching storage diffs: %w", fetchErr)
		case diff := <-diffsChan:
			extractor.persistDiff(diff)
		}
	}
}

func (extractor DiffExtractor) persistDiff(rawDiff types.RawDiff) {
	_, err := extractor.StorageDiffRepository.CreateStorageDiff(rawDiff)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Tracef("ignoring duplicate diff. Block number: %v, blockHash: %v, storageKey: %v, storageValue: %v",
				rawDiff.BlockHeight, rawDiff.BlockHash.Hex(), rawDiff.StorageKey, rawDiff.StorageValue)
			return
		}
		logrus.Warnf("failed to persist storage diff: %s", err.Error())
	}
}
