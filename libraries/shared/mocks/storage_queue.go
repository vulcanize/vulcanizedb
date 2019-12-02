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
	AddPassedDiffs  map[int]utils.StorageDiff
	DeleteErr       error
	DeletePassedIds []int
	GetAllErr       error
	DiffsToReturn   map[int]utils.StorageDiff
	GetAllCalled    bool
}

// Add mock method
func (queue *MockStorageQueue) Add(diff utils.StorageDiff) error {
	queue.AddCalled = true
	if queue.AddPassedDiffs == nil {
		queue.AddPassedDiffs = make(map[int]utils.StorageDiff)
	}
	queue.AddPassedDiffs[diff.ID] = diff
	return queue.AddError
}

// Delete mock method
func (queue *MockStorageQueue) Delete(id int) error {
	queue.DeletePassedIds = append(queue.DeletePassedIds, id)
	delete(queue.DiffsToReturn, id)
	return queue.DeleteErr
}

// GetAll mock method
func (queue *MockStorageQueue) GetAll() ([]utils.StorageDiff, error) {
	queue.GetAllCalled = true
	diffs := make([]utils.StorageDiff, 0)
	for _, diff := range queue.DiffsToReturn {
		diffs = append(diffs, diff)
	}
	return diffs, queue.GetAllErr
}
