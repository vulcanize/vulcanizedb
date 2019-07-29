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

type MockStorageQueue struct {
	AddCalled      bool
	AddError       error
	AddPassedDiff  utils.StorageDiff
	DeleteErr      error
	DeletePassedId int
	GetAllErr      error
	DiffsToReturn  []utils.StorageDiff
}

func (queue *MockStorageQueue) Add(diff utils.StorageDiff) error {
	queue.AddCalled = true
	queue.AddPassedDiff = diff
	return queue.AddError
}

func (queue *MockStorageQueue) Delete(id int) error {
	queue.DeletePassedId = id
	return queue.DeleteErr
}

func (queue *MockStorageQueue) GetAll() ([]utils.StorageDiff, error) {
	return queue.DiffsToReturn, queue.GetAllErr
}
