package vat_tune

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type VatTuneRepository struct {
	db *postgres.DB
}

func (repository VatTuneRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		vatTune, ok := model.(VatTuneModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VatTuneModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.vat_tune (header_id, ilk, urn, v, w, dink, dart, tx_idx, log_idx, raw_log)
	   VALUES($1, $2, $3, $4, $5, $6::NUMERIC, $7::NUMERIC, $8, $9, $10)`,
			headerID, vatTune.Ilk, vatTune.Urn, vatTune.V, vatTune.W, vatTune.Dink, vatTune.Dart, vatTune.TransactionIndex, vatTune.LogIndex, vatTune.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatTuneChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository VatTuneRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.VatTuneChecked)
}

func (repository *VatTuneRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
