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
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type MockStorageRepository struct {
	CreateErr      error
	PassedHeaderID int64
	PassedDiffID   int64
	PassedMetadata utils.StorageValueMetadata
	PassedValue    interface{}
}

func (repository *MockStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	repository.PassedHeaderID = headerID
	repository.PassedDiffID = diffID
	repository.PassedMetadata = metadata
	repository.PassedValue = value
	return repository.CreateErr
}

func (*MockStorageRepository) SetDB(db *postgres.DB) {
	panic("implement me")
}
