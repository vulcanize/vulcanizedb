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

package mocks

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

type MockStorageDiffRepository struct {
	CreatePassedRawDiffs   []types.RawDiff
	CreateReturnID         int64
	CreateReturnError      error
	GetNewDiffsDiffs       []types.PersistedDiff
	GetNewDiffsErrors      []error
	MarkCheckedPassedID    int64
	MarkFromBackfillCalled bool
	MarkFromBackfillError  error
}

func (repository *MockStorageDiffRepository) GetNewDiffs(diffs chan types.PersistedDiff, errs chan error, done chan bool) {
	for _, diff := range repository.GetNewDiffsDiffs {
		diffs <- diff
	}
	for _, err := range repository.GetNewDiffsErrors {
		errs <- err
	}
}

func (repository *MockStorageDiffRepository) MarkChecked(id int64) error {
	repository.MarkCheckedPassedID = id
	return nil
}

func (repository *MockStorageDiffRepository) CreateStorageDiff(rawDiff types.RawDiff) (int64, error) {
	repository.CreatePassedRawDiffs = append(repository.CreatePassedRawDiffs, rawDiff)
	return repository.CreateReturnID, repository.CreateReturnError
}

func (repository *MockStorageDiffRepository) MarkFromBackfill(id int64) error {
	repository.MarkFromBackfillCalled = true
	return repository.MarkFromBackfillError
}
