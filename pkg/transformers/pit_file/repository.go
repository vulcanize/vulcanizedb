package pit_file

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, transactionIndex uint, model PitFileModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type PitFileRepository struct {
	db *postgres.DB
}

func NewPitFileRepository(db *postgres.DB) PitFileRepository {
	return PitFileRepository{
		db: db,
	}
}

func (repository PitFileRepository) Create(headerID int64, transactionIndex uint, model PitFileModel) error {
	_, err := repository.db.Exec(
		`INSERT into maker.pit_file (header_id, ilk, what, risk, tx_idx, raw_log)
        VALUES($1, $2, $3, $4::NUMERIC, $5, $6)`,
		headerID, model.Ilk, model.What, model.Risk, model.TransactionIndex, model.Raw,
	)
	return err
}

func (repository PitFileRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.pit_file on headers.id = header_id
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
