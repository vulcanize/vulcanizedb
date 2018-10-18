// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ilk

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type DripFileIlkRepository struct {
	DB *postgres.DB
}

func (repository DripFileIlkRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.DB.Begin()
	if err != nil {
		return err
	}

	var ilk DripFileIlkModel
	for _, model := range models {
		ilk = model.(DripFileIlkModel)
		_, err = tx.Exec(
			`INSERT into maker.drip_file_ilk (header_id, ilk, vow, tax, log_idx, tx_idx, raw_log)
        	VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)`,
			headerID, ilk.Ilk, ilk.Vow, ilk.Tax, ilk.LogIndex, ilk.TransactionIndex, ilk.Raw,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, drip_file_ilk_checked)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET drip_file_ilk_checked = $2`, headerID, true)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repository DripFileIlkRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.DB.Exec(`INSERT INTO public.checked_headers (header_id, drip_file_ilk_checked)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET drip_file_ilk_checked = $2`, headerID, true)
	return err
}

func (repository DripFileIlkRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.DB.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               	LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR drip_file_ilk_checked IS FALSE)
               	AND headers.block_number >= $1
               	AND headers.block_number <= $2
               	AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.DB.Node.ID,
	)
	return result, err
}

func (repository *DripFileIlkRepository) SetDB(db *postgres.DB) {
	repository.DB = db
}
