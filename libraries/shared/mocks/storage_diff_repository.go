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
	CreatePassedRawDiffs []types.RawDiff
	CreateReturnID       int64
	CreateReturnError    error
}

func (repository *MockStorageDiffRepository) CreateStorageDiff(rawDiff types.RawDiff) (int64, error) {
	repository.CreatePassedRawDiffs = append(repository.CreatePassedRawDiffs, rawDiff)
	return repository.CreateReturnID, repository.CreateReturnError
}
