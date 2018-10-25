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

package flog

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type FlogRepository struct {
	db *postgres.DB
}

func (repository FlogRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		flog, ok := model.(FlogModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, FlogModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.vow_flog (header_id, era, log_idx, tx_idx, raw_log)
			VALUES($1, $2::NUMERIC, $3, $4, $5)`,
			headerID, flog.Era, flog.LogIndex, flog.TransactionIndex, flog.Raw,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, vow_flog_checked)
		VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET vow_flog_checked = $2`, headerID, true)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository FlogRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, vow_flog_checked)
		VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET vow_flog_checked = $2`, headerID, true)
	return err
}

func (repository FlogRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
			LEFT JOIN checked_headers on headers.id = header_id
		WHERE (header_id ISNULL OR vow_flog_checked IS FALSE)
			AND headers.block_number >= $1
			AND headers.block_number <= $2
			AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}

func (repository *FlogRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
