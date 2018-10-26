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

package dent

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type DentRepository struct {
	db *postgres.DB
}

func (r DentRepository) Create(headerId int64, models []interface{}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		dent, ok := model.(DentModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, DentModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.dent (header_id, bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			headerId, dent.BidId, dent.Lot, dent.Bid, dent.Guy, dent.Tic, dent.LogIndex, dent.TransactionIndex, dent.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, dent_checked)
			VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET dent_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r DentRepository) MarkHeaderChecked(headerId int64) error {
	_, err := r.db.Exec(`INSERT INTO public.checked_headers (header_id, dent_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET dent_checked = $2`, headerId, true)
	return err
}

func (r DentRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var missingHeaders []core.Header

	err := r.db.Select(
		&missingHeaders,
		`SELECT headers.id, headers.block_number FROM headers
			LEFT JOIN checked_headers on headers.id = header_id
		WHERE (header_id ISNULL OR dent_checked IS FALSE)
			AND headers.block_number >= $1
			AND headers.block_number <= $2
			AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		r.db.Node.ID,
	)

	return missingHeaders, err
}

func (repository *DentRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
