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
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type DentRepository struct {
	db *postgres.DB
}

func (repository DentRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		dent, ok := model.(DentModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, DentModel{})
		}

		err = shared.ValidateHeaderConsistency(headerID, dent.Raw, repository.db)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			`INSERT into maker.dent (header_id, bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			headerID, dent.BidId, dent.Lot, dent.Bid, dent.Guy, dent.Tic, dent.LogIndex, dent.TransactionIndex, dent.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.DentChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository DentRepository) MarkHeaderChecked(headerId int64) error {
	return shared.MarkHeaderChecked(headerId, repository.db, constants.DentChecked)
}

func (repository DentRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.DentChecked)
}

func (repository *DentRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
