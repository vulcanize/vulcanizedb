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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerId int64, models []VatFluxModel) error
	MissingHeaders(startingBlock, endingBlock int64) ([]core.Header, error)
	MarkCheckedHeader(headerId int64) error
}

type VatFluxRepository struct {
	DB *postgres.DB
}

func NewVatFluxRepository(db *postgres.DB) VatFluxRepository {
	return VatFluxRepository{DB: db}
}

func (repository VatFluxRepository) Create(headerId int64, models []VatFluxModel) error {
	tx, err := repository.DB.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		_, err := tx.Exec(`INSERT INTO maker.vat_flux (header_id, ilk, dst, src, rad, tx_idx, log_idx, raw_log)
		VALUES($1, $2, $3, $4, $5::NUMERIC, $6, $7, $8)`,
			headerId, model.Ilk, model.Dst, model.Src, model.Rad, model.TransactionIndex, model.LogIndex, model.Raw)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, vat_flux_checked)
			VALUES($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET vat_flux_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatFluxRepository) MissingHeaders(startingBlock, endingBlock int64) ([]core.Header, error) {
	var headers []core.Header
	err := repository.DB.Select(&headers,
		`SELECT headers.id, block_number from headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR vat_flux_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlock, endingBlock, repository.DB.Node.ID)

	return headers, err
}

func (repository VatFluxRepository) MarkCheckedHeader(headerId int64) error {
	_, err := repository.DB.Exec(`INSERT INTO public.checked_headers (header_id, vat_flux_checked)
			VALUES($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET vat_flux_checked = $2`, headerId, true)

	return err
}
