package vat_grab

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type VatGrabRepository struct {
	db *postgres.DB
}

func (repository VatGrabRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		vatGrab, ok := model.(VatGrabModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, VatGrabModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.vat_grab (header_id, ilk, urn, v, w, dink, dart, log_idx, tx_idx, raw_log)
	   VALUES($1, $2, $3, $4, $5, $6::NUMERIC, $7::NUMERIC, $8, $9, $10)`,
			headerID, vatGrab.Ilk, vatGrab.Urn, vatGrab.V, vatGrab.W, vatGrab.Dink, vatGrab.Dart, vatGrab.LogIndex, vatGrab.TransactionIndex, vatGrab.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatGrabChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatGrabRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.VatGrabChecked)
}

func (repository *VatGrabRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
