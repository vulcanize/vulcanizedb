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

package storage

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Queue interface {
	Add(diff types.PersistedDiff) error
	Delete(id int64) error
	GetAll() ([]types.PersistedDiff, error)
}

type queue struct {
	db *postgres.DB
}

func NewStorageQueue(db *postgres.DB) queue {
	return queue{db: db}
}

func (queue queue) Add(diff types.PersistedDiff) error {
	_, err := queue.db.Exec(`INSERT INTO public.queued_storage (diff_id) VALUES
		($1) ON CONFLICT DO NOTHING`, diff.ID)
	return err
}

func (queue queue) Delete(diffID int64) error {
	_, err := queue.db.Exec(`DELETE FROM public.queued_storage WHERE diff_id = $1`, diffID)
	return err
}

func (queue queue) GetAll() ([]types.PersistedDiff, error) {
	var result []types.PersistedDiff
	err := queue.db.Select(&result, `SELECT storage_diff.id, hashed_address, block_height, block_hash, storage_key, storage_value
		FROM public.queued_storage
			LEFT JOIN public.storage_diff ON queued_storage.diff_id = storage_diff.id`)
	return result, err
}
