package debt_ceiling

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, model PitFileDebtCeilingModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type PitFileDebtCeilingRepository struct {
	db *postgres.DB
}

func NewPitFileDebtCeilingRepository(db *postgres.DB) PitFileDebtCeilingRepository {
	return PitFileDebtCeilingRepository{
		db: db,
	}
}

func (repository PitFileDebtCeilingRepository) Create(headerID int64, model PitFileDebtCeilingModel) error {
	_, err := repository.db.Exec(
		`INSERT into maker.pit_file_debt_ceiling (header_id, what, data, tx_idx, raw_log)
        VALUES($1, $2, $3::NUMERIC, $4, $5)`,
		headerID, model.What, model.Data, model.TransactionIndex, model.Raw,
	)
	return err
}

func (repository PitFileDebtCeilingRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.pit_file_debt_ceiling on headers.id = header_id
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
