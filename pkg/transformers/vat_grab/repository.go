package vat_grab

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, models []VatGrabModel) error
	MarkHeaderChecked(headerID int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type VatGrabRepository struct {
	db *postgres.DB
}

func NewVatGrabRepository(db *postgres.DB) VatGrabRepository {
	return VatGrabRepository{
		db: db,
	}
}

func (repository VatGrabRepository) Create(headerID int64, models []VatGrabModel) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	for _, model := range models {
		_, err = tx.Exec(
			`INSERT into maker.vat_grab (header_id, ilk, urn, v, w, dink, dart, tx_idx, raw_log)
	   VALUES($1, $2, $3, $4, $5, $6::NUMERIC, $7::NUMERIC, $8, $9)`,
			headerID, model.Ilk, model.Urn, model.V, model.W, model.Dink, model.Dart, model.TransactionIndex, model.Raw,
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
