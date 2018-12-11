package vat_slip

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type VatSlipRepository struct {
	db *postgres.DB
}

func (repository VatSlipRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		vatSlip, ok := model.(VatSlipModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, VatSlipModel{})
		}

		_, err = tx.Exec(
			`INSERT into maker.vat_slip (header_id, ilk, guy, rad, tx_idx, log_idx, raw_log)
			VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)`,
			headerID, vatSlip.Ilk, vatSlip.Guy, vatSlip.Rad, vatSlip.TransactionIndex, vatSlip.LogIndex, vatSlip.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.Repository{}.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatSlipChecked)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repository VatSlipRepository) MarkHeaderChecked(headerID int64) error {
	return shared.Repository{}.MarkHeaderChecked(headerID, repository.db, constants.VatSlipChecked)
}

func (repository *VatSlipRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
