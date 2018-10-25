package vat_slip

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
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
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, vat_slip_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET vat_slip_checked = $2`, headerID, true)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repository VatSlipRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, vat_slip_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET vat_slip_checked = $2`, headerID, true)
	return err
}

func (repository VatSlipRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR vat_slip_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}

func (repository *VatSlipRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
