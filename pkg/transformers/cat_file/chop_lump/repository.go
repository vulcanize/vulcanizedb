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

package chop_lump

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type CatFileChopLumpRepository struct {
	db *postgres.DB
}

func (repository CatFileChopLumpRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}

	for _, model := range models {
		chopLump, ok := model.(CatFileChopLumpModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, CatFileChopLumpModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(chopLump.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		_, execErr := tx.Exec(
			`INSERT into maker.cat_file_chop_lump (header_id, ilk, what, data, tx_idx, log_idx, raw_log)
			VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)
			ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET ilk = $2, what = $3, data = $4, raw_log = $7;`,
			headerID, ilkID, chopLump.What, chopLump.Data, chopLump.TransactionIndex, chopLump.LogIndex, chopLump.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.CatFileChopLumpChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository CatFileChopLumpRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.CatFileChopLumpChecked)
}

func (repository CatFileChopLumpRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.CatFileChopLumpChecked)
}

func (repository CatFileChopLumpRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.CatFileChopLumpChecked)
}

func (repository *CatFileChopLumpRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
