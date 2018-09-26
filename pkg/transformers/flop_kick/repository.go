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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerId int64, flopKick []Model) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
	MarkHeaderChecked(headerId int64) error
}

type FlopKickRepository struct {
	DB *postgres.DB
}

func NewFlopKickRepository(db *postgres.DB) FlopKickRepository {
	return FlopKickRepository{DB: db}
}

func (r FlopKickRepository) Create(headerId int64, flopKicks []Model) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	for _, flopKick := range flopKicks {
		_, err = tx.Exec(
			`INSERT into maker.flop_kick (header_id, bid_id, lot, bid, gal, "end", tx_idx, raw_log)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
			headerId, flopKick.BidId, flopKick.Lot, flopKick.Bid, flopKick.Gal, flopKick.End, flopKick.TransactionIndex, flopKick.Raw,
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
	_, err := r.DB.Exec(`INSERT INTO public.checked_headers (header_id, flop_kick_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET flop_kick_checked = $2`, headerId, true)
	return err
}

func (r FlopKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := r.DB.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR flop_kick_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		r.DB.Node.ID,
	)

	return result, err
}
