package ilk

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, model PitFileIlkModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type PitFileIlkRepository struct {
	db *postgres.DB
}

func NewPitFileIlkRepository(db *postgres.DB) PitFileIlkRepository {
	return PitFileIlkRepository{
		db: db,
	}
}

func (repository PitFileIlkRepository) Create(headerID int64, model PitFileIlkModel) error {
	_, err := repository.db.Exec(
		`INSERT into maker.pit_file_ilk (header_id, ilk, what, data, tx_idx, raw_log)
        VALUES($1, $2, $3, $4::NUMERIC, $5, $6)`,
		headerID, model.Ilk, model.What, model.Data, model.TransactionIndex, model.Raw,
	)
	return err
}

func (repository PitFileIlkRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.pit_file_ilk on headers.id = header_id
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
