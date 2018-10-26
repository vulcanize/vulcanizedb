package vat_grab

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, vat_grab_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET vat_grab_checked = $2`, headerID, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository VatGrabRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, vat_grab_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET vat_grab_checked = $2`, headerID, true)
	return err
}

func (repository VatGrabRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
	          LEFT JOIN checked_headers on headers.id = header_id
	          WHERE (header_id ISNULL OR vat_grab_checked IS FALSE)
	          AND headers.block_number >= $1
	          AND headers.block_number <= $2
	          AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}

func (repository *VatGrabRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
