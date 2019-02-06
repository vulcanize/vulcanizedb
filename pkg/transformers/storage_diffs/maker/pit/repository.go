// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package pit

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type PitStorageRepository struct {
	db *postgres.DB
}

func (repository *PitStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository PitStorageRepository) Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case IlkLine:
		return repository.insertIlkLine(blockNumber, blockHash, metadata.Keys[shared.Ilk], value.(string))
	case IlkSpot:
		return repository.insertIlkSpot(blockNumber, blockHash, metadata.Keys[shared.Ilk], value.(string))
	case PitDrip:
		return repository.insertPitDrip(blockNumber, blockHash, value.(string))
	case PitLine:
		return repository.insertPitLine(blockNumber, blockHash, value.(string))
	case PitLive:
		return repository.insertPitLive(blockNumber, blockHash, value.(string))
	case PitVat:
		return repository.insertPitVat(blockNumber, blockHash, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized pit contract storage name: %s", metadata.Name))
	}
}

func (repository PitStorageRepository) insertIlkLine(blockNumber int, blockHash string, ilk string, line string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pit_ilk_line (block_number, block_hash, ilk, line) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilk, line)
	return err
}

func (repository PitStorageRepository) insertIlkSpot(blockNumber int, blockHash string, ilk string, spot string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pit_ilk_spot (block_number, block_hash, ilk, spot) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilk, spot)
	return err
}

func (repository PitStorageRepository) insertPitDrip(blockNumber int, blockHash string, drip string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pit_drip (block_number, block_hash, drip) VALUES ($1, $2, $3)`, blockNumber, blockHash, drip)
	return err
}

func (repository PitStorageRepository) insertPitLine(blockNumber int, blockHash string, line string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pit_line (block_number, block_hash, line) VALUES ($1, $2, $3)`, blockNumber, blockHash, line)
	return err
}

func (repository PitStorageRepository) insertPitLive(blockNumber int, blockHash string, live string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pit_live (block_number, block_hash, live) VALUES ($1, $2, $3)`, blockNumber, blockHash, live)
	return err
}

func (repository PitStorageRepository) insertPitVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pit_vat (block_number, block_hash, vat) VALUES ($1, $2, $3)`, blockNumber, blockHash, vat)
	return err
}
