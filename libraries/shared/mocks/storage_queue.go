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
	AddPassedDiffs  []utils.StorageDiff
	DeleteErr       error
	DeletePassedIds []int
	GetAllErr       error
	DiffsToReturn   []utils.StorageDiff
}

// Add mock method
func (queue *MockStorageQueue) Add(diff utils.StorageDiff) error {
	queue.AddCalled = true
	queue.AddPassedDiffs = append(queue.AddPassedDiffs, diff)
	return queue.AddError
}

// Delete mock method
func (queue *MockStorageQueue) Delete(id int) error {
	queue.DeletePassedIds = append(queue.DeletePassedIds, id)
	return queue.DeleteErr
}

// GetAll mock method
func (queue *MockStorageQueue) GetAll() ([]utils.StorageDiff, error) {
	return queue.DiffsToReturn, queue.GetAllErr
}
