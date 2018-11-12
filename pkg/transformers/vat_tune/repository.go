package vat_tune

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
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

		err = shared.ValidateHeaderConsistency(headerID, vatTune.Raw, repository.db)
		if err != nil {
			tx.Rollback()
			return err
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

	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, vat_tune_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET vat_tune_checked = $2`, headerID, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatTuneRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, vat_tune_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET vat_tune_checked = $2`, headerID, true)
	return err
}

func (repository VatTuneRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
	          LEFT JOIN checked_headers on headers.id = header_id
	          WHERE (header_id ISNULL OR vat_tune_checked IS FALSE)
	          AND headers.block_number >= $1
	          AND headers.block_number <= $2
	          AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}

func (repository *VatTuneRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
