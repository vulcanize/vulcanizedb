package vat_tune

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type VatTuneRepository struct {
	db *postgres.DB
}

func (repository VatTuneRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		vatTune, ok := model.(VatTuneModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, VatTuneModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.vat_tune (header_id, ilk, urn, v, w, dink, dart, tx_idx, log_idx, raw_log)
	   VALUES($1, $2, $3, $4, $5, $6::NUMERIC, $7::NUMERIC, $8, $9, $10)`,
			headerID, vatTune.Ilk, vatTune.Urn, vatTune.V, vatTune.W, vatTune.Dink, vatTune.Dart, vatTune.TransactionIndex, vatTune.LogIndex, vatTune.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.Repository{}.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatTuneChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatTuneRepository) MarkHeaderChecked(headerID int64) error {
	return shared.Repository{}.MarkHeaderChecked(headerID, repository.db, constants.VatTuneChecked)
}

func (repository *VatTuneRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
