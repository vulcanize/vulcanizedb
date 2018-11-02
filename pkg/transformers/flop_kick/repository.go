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

package flop_kick

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type FlopKickRepository struct {
	db *postgres.DB
}

func (r FlopKickRepository) Create(headerId int64, models []interface{}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	for _, flopKick := range models {
		flopKickModel, ok := flopKick.(Model)

		if !ok {
			return fmt.Errorf("model of type %T, not %T", flopKick, Model{})
		}
		_, err = tx.Exec(
			`INSERT into maker.flop_kick (header_id, bid_id, lot, bid, gal, "end", tx_idx, log_idx, raw_log)
        VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6, $7, $8, $9)`,
			headerId, flopKickModel.BidId, flopKickModel.Lot, flopKickModel.Bid, flopKickModel.Gal, flopKickModel.End, flopKickModel.TransactionIndex, flopKickModel.LogIndex, flopKickModel.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, flop_kick_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET flop_kick_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r FlopKickRepository) MarkHeaderChecked(headerId int64) error {
	_, err := r.db.Exec(`INSERT INTO public.checked_headers (header_id, flop_kick_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET flop_kick_checked = $2`, headerId, true)
	return err
}

func (r FlopKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := r.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR flop_kick_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		r.db.Node.ID,
	)

	return result, err
}

func (repository *FlopKickRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
