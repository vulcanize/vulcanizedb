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
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"log"
)

type DentRepository struct {
	db *postgres.DB
}

func (repository DentRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}

	tic, gettTicErr := shared.GetTicInTx(headerID, tx)
	if gettTicErr != nil {
		return gettTicErr
	}

	for _, model := range models {
		dent, ok := model.(DentModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, DentModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.dent (header_id, bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			headerID, dent.BidId, dent.Lot, dent.Bid, dent.Guy, tic, dent.LogIndex, dent.TransactionIndex, dent.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.DentChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Println("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository DentRepository) MarkHeaderChecked(headerId int64) error {
	return shared.MarkHeaderChecked(headerId, repository.db, constants.DentChecked)
}

func (repository *DentRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
