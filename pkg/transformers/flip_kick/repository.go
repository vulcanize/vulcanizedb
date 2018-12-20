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

package flip_kick

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type FlipKickRepository struct {
	db *postgres.DB
}

func (repository FlipKickRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		flipKickModel, ok := model.(FlipKickModel)
		if !ok {
			return fmt.Errorf("model of type %T, not %T", model, FlipKickModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.flip_kick (header_id, bid_id, lot, bid, gal, "end", urn, tab, tx_idx, log_idx, raw_log)
        VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6, $7, $8::NUMERIC, $9, $10, $11)`,
			headerID, flipKickModel.BidId, flipKickModel.Lot, flipKickModel.Bid, flipKickModel.Gal, flipKickModel.End, flipKickModel.Urn, flipKickModel.Tab, flipKickModel.TransactionIndex, flipKickModel.LogIndex, flipKickModel.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}
	checkHeaderErr := shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.FlipKickChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository FlipKickRepository) MarkHeaderChecked(headerId int64) error {
	return shared.MarkHeaderChecked(headerId, repository.db, constants.FlipKickChecked)
}

func (repository *FlipKickRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
