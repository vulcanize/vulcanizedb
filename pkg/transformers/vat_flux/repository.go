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

package vat_flux

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type VatFluxRepository struct {
	db *postgres.DB
}

func (repository VatFluxRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		vatFlux, ok := model.(VatFluxModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, VatFluxModel{})
		}

		_, err = tx.Exec(`INSERT INTO maker.vat_flux (header_id, ilk, dst, src, rad, tx_idx, log_idx, raw_log)
		VALUES($1, $2, $3, $4, $5::NUMERIC, $6, $7, $8)`,
			headerID, vatFlux.Ilk, vatFlux.Dst, vatFlux.Src, vatFlux.Rad, vatFlux.TransactionIndex, vatFlux.LogIndex, vatFlux.Raw)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatFluxChecked)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repository VatFluxRepository) MarkHeaderChecked(headerId int64) error {
	return shared.MarkHeaderChecked(headerId, repository.db, constants.VatFluxChecked)
}

func (repository *VatFluxRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
