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
	case VowCow:
		return repository.insertVowCow(blockNumber, blockHash, value.(string))
	case VowRow:
		return repository.insertVowRow(blockNumber, blockHash, value.(string))
	case VowSin:
		return repository.insertVowSin(blockNumber, blockHash, value.(string))
	case VowAsh:
		return repository.insertVowAsh(blockNumber, blockHash, value.(string))
	case VowWait:
		return repository.insertVowWait(blockNumber, blockHash, value.(string))
	case VowSump:
		return repository.insertVowSump(blockNumber, blockHash, value.(string))
	case VowBump:
		return repository.insertVowBump(blockNumber, blockHash, value.(string))
	case VowHump:
		return repository.insertVowHump(blockNumber, blockHash, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository VowStorageRepository) insertVowVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_vat (block_number, block_hash, vat) VALUES ($1, $2, $3)`, blockNumber, blockHash, vat)

	return err
}

func (repository VowStorageRepository) insertVowCow(blockNumber int, blockHash string, cow string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_cow (block_number, block_hash, cow) VALUES ($1, $2, $3)`, blockNumber, blockHash, cow)

	return err
}

func (repository VowStorageRepository) insertVowRow(blockNumber int, blockHash string, row string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_row (block_number, block_hash, row) VALUES ($1, $2, $3)`, blockNumber, blockHash, row)

	return err
}

func (repository VowStorageRepository) insertVowSin(blockNumber int, blockHash string, sin string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_sin (block_number, block_hash, sin) VALUES ($1, $2, $3)`, blockNumber, blockHash, sin)

	return err
}

func (repository VowStorageRepository) insertVowAsh(blockNumber int, blockHash string, ash string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_ash (block_number, block_hash, ash) VALUES ($1, $2, $3)`, blockNumber, blockHash, ash)

	return err
}

func (repository VowStorageRepository) insertVowWait(blockNumber int, blockHash string, wait string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_wait (block_number, block_hash, wait) VALUES ($1, $2, $3)`, blockNumber, blockHash, wait)

	return err
}

func (repository VowStorageRepository) insertVowSump(blockNumber int, blockHash string, sump string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_sump (block_number, block_hash, sump) VALUES ($1, $2, $3)`, blockNumber, blockHash, sump)

	return err
}

func (repository VowStorageRepository) insertVowBump(blockNumber int, blockHash string, bump string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_bump (block_number, block_hash, bump) VALUES ($1, $2, $3)`, blockNumber, blockHash, bump)

	return err
}

func (repository VowStorageRepository) insertVowHump(blockNumber int, blockHash string, hump string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vow_hump (block_number, block_hash, hump) VALUES ($1, $2, $3)`, blockNumber, blockHash, hump)

	return err
}
