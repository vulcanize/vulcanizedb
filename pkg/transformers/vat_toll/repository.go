package vat_toll

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type VatTollRepository struct {
	db *postgres.DB
}

func (repository VatTollRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		vatToll, ok := model.(VatTollModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, VatTollModel{})
		}

		err = shared.ValidateHeaderConsistency(headerID, vatToll.Raw, repository.db)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			`INSERT into maker.vat_toll (header_id, ilk, urn, take, tx_idx, log_idx, raw_log)
			VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)`,
			headerID, vatToll.Ilk, vatToll.Urn, vatToll.Take, vatToll.TransactionIndex, vatToll.LogIndex, vatToll.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = shared.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatTollChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatTollRepository) MarkHeaderChecked(headerID int64) error {
	return shared.MarkHeaderChecked(headerID, repository.db, constants.VatTollChecked)
}

func (repository VatTollRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return shared.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.VatTollChecked)
}

func (repository *VatTollRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
