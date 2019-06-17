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
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type IStorageQueue interface {
	Add(row utils.StorageDiffRow) error
	Delete(id int) error
	GetAll() ([]utils.StorageDiffRow, error)
}

type StorageQueue struct {
	db *postgres.DB
}

func NewStorageQueue(db *postgres.DB) StorageQueue {
	return StorageQueue{db: db}
}

func (queue StorageQueue) Add(row utils.StorageDiffRow) error {
	_, err := queue.db.Exec(`INSERT INTO public.queued_storage (contract,
		block_hash, block_height, storage_key, storage_value) VALUES
		($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`, row.Contract.Bytes(), row.BlockHash.Bytes(),
		row.BlockHeight, row.StorageKey.Bytes(), row.StorageValue.Bytes())
	return err
}

func (queue StorageQueue) Delete(id int) error {
	_, err := queue.db.Exec(`DELETE FROM public.queued_storage WHERE id = $1`, id)
	return err
}

func (queue StorageQueue) GetAll() ([]utils.StorageDiffRow, error) {
	var result []utils.StorageDiffRow
	err := queue.db.Select(&result, `SELECT * FROM public.queued_storage`)
	return result, err
}
