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

package frob

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type FrobRepository struct {
	db *postgres.DB
}

func (repository FrobRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		frobModel, ok := model.(FrobModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, FrobModel{})
		}
		_, err = tx.Exec(`INSERT INTO maker.frob (header_id, art, dart, dink, iart, ilk, ink, urn, raw_log, log_idx, tx_idx)
		VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7::NUMERIC, $8, $9, $10, $11)`,
			headerID, frobModel.Art, frobModel.Dart, frobModel.Dink, frobModel.IArt, frobModel.Ilk, frobModel.Ink, frobModel.Urn, frobModel.Raw, frobModel.LogIndex, frobModel.TransactionIndex)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.FrobChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository FrobRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.FrobChecked)
}

func (repository FrobRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.FrobChecked)
}

func (repository *FrobRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
