package vat_init

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, model VatInitModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type VatInitRepository struct {
	db *postgres.DB
}

func NewVatInitRepository(db *postgres.DB) VatInitRepository {
	return VatInitRepository{
		db: db,
	}
}

func (repository VatInitRepository) Create(headerID int64, model VatInitModel) error {
	_, err := repository.db.Exec(
		`INSERT into maker.vat_init (header_id, ilk, tx_idx, raw_log)
        VALUES($1, $2, $3, $4)`,
		headerID, model.Ilk, model.TransactionIndex, model.Raw,
	)
	return err
}

func (repository VatInitRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.vat_init on headers.id = header_id
               WHERE header_id ISNULL
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)

	return result, err
}
