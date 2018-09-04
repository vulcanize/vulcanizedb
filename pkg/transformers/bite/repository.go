package bite

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, model BiteModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type BiteRepository struct {
	db *postgres.DB
}

func NewBiteRepository(db *postgres.DB) Repository {
	return BiteRepository{db: db}
}

func (repository BiteRepository) Create(headerID int64, model BiteModel) error {
	_, err := repository.db.Exec(
		`INSERT into maker.bite (header_id, id, ilk, lad, ink, art, iart, tab, flip, tx_idx, raw_log)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		headerID, model.Id, model.Ilk, model.Lad, model.Ink, model.Art, model.IArt, model.Tab, model.Flip, model.TransactionIndex, model.Raw,
	)

	if err != nil {
		return err
	}

	return nil
}
func (repository BiteRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.bite on headers.id = header_id
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
