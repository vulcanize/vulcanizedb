// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package flap_kick

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type FlapKickRepository struct {
	db *postgres.DB
}

func (repository *FlapKickRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		flapKickModel, ok := model.(FlapKickModel)
		if !ok {
			return fmt.Errorf("model of type %T, not %T", model, FlapKickModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.flap_kick (header_id, bid_id, lot, bid, gal, "end", tx_idx, log_idx, raw_log)
        VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6, $7, $8, $9)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET bid_id = $2, lot = $3, bid = $4, gal = $5, "end" = $6, raw_log = $9;`,
			headerID, flapKickModel.BidId, flapKickModel.Lot, flapKickModel.Bid, flapKickModel.Gal, flapKickModel.End, flapKickModel.TransactionIndex, flapKickModel.LogIndex, flapKickModel.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.FlapKickChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository *FlapKickRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.FlapKickChecked)
}

func (repository FlapKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.FlapKickChecked)
}

func (repository FlapKickRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.FlapKickChecked)
}

func (repository *FlapKickRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
