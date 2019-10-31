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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

// MockStorageQueue for tests
type MockStorageQueue struct {
	AddCalled       bool
	AddError        error
	AddPassedDiffs  []utils.PersistedStorageDiff
	DeleteErr       error
	DeletePassedIds []int64
	GetAllErr       error
	DiffsToReturn   []utils.PersistedStorageDiff
	GetAllCalled    bool
}

// Add mock method
func (queue *MockStorageQueue) Add(diff utils.PersistedStorageDiff) error {
	queue.AddCalled = true
	queue.AddPassedDiffs = append(queue.AddPassedDiffs, diff)
	return queue.AddError
}

// Delete mock method
func (queue *MockStorageQueue) Delete(id int64) error {
	queue.DeletePassedIds = append(queue.DeletePassedIds, id)
	var diffs []utils.PersistedStorageDiff
	for _, diff := range queue.DiffsToReturn {
		if diff.ID != id {
			diffs = append(diffs, diff)
		}
	}
	queue.DiffsToReturn = diffs
	return queue.DeleteErr
}

// GetAll mock method
func (queue *MockStorageQueue) GetAll() ([]utils.PersistedStorageDiff, error) {
	queue.GetAllCalled = true
	return queue.DiffsToReturn, queue.GetAllErr
}
