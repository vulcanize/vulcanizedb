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

package deal

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerId int64, models []DealModel) error
	MarkHeaderChecked(headerID int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}
type DealRepository struct {
	db *postgres.DB
}

func NewDealRepository(database *postgres.DB) DealRepository {
	return DealRepository{db: database}
}
func (r DealRepository) Create(headerId int64, models []DealModel) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		_, err = tx.Exec(
			`INSERT into maker.deal (header_id, bid_id, tx_idx, raw_log)
         VALUES($1, $2, $3, $4)`,
			headerId, model.BidId, model.TransactionIndex, model.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, deal_checked)
			VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET deal_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r DealRepository) MarkHeaderChecked(headerID int64) error {
	_, err := r.db.Exec(`INSERT INTO public.checked_headers (header_id, deal_checked)
		VALUES ($1, $2) 
	ON CONFLICT (header_id) DO
		UPDATE SET deal_checked = $2`, headerID, true)
	return err
}

func (r DealRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var missingHeaders []core.Header
	err := r.db.Select(&missingHeaders,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR deal_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		r.db.Node.ID,
	)
	return missingHeaders, err
}
