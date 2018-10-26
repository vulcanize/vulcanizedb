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

package vat_heal

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type VatHealRepository struct {
	db *postgres.DB
}

func (repository *VatHealRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VatHealRepository) Create(headerId int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		vatHeal, ok := model.(VatHealModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, VatHealModel{})
		}
		_, err := tx.Exec(`INSERT INTO maker.vat_heal (header_id, urn, v, rad, log_idx, tx_idx, raw_log)
		VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)`,
			headerId, vatHeal.Urn, vatHeal.V, vatHeal.Rad, vatHeal.LogIndex, vatHeal.TransactionIndex, vatHeal.Raw)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, vat_heal_checked)
			VALUES($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET vat_heal_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatHealRepository) MissingHeaders(startingBlock, endingBlock int64) ([]core.Header, error) {
	var headers []core.Header
	err := repository.db.Select(&headers,
		`SELECT headers.id, block_number from headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR vat_heal_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlock, endingBlock, repository.db.Node.ID)

	return headers, err
}

func (repository VatHealRepository) MarkHeaderChecked(headerId int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, vat_heal_checked)
			VALUES($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET vat_heal_checked = $2`, headerId, true)

	return err
}
