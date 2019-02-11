/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package vow

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type VowStorageRepository struct {
	db *postgres.DB
}

func (repository *VowStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VowStorageRepository) Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case VowVat:
		return repository.insertVowVat(blockNumber, blockHash, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository VowStorageRepository) insertVowVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_vat (block_number, block_hash, vat) VALUES ($1, $2, $3)`, blockNumber, blockHash, vat)

	return err
}
