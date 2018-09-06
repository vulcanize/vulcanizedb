package stability_fee

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerID int64, model PitFileStabilityFeeModel) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type PitFileStabilityFeeRepository struct {
	db *postgres.DB
}

func NewPitFileStabilityFeeRepository(db *postgres.DB) PitFileStabilityFeeRepository {
	return PitFileStabilityFeeRepository{
		db: db,
	}
}

func (repository PitFileStabilityFeeRepository) Create(headerID int64, model PitFileStabilityFeeModel) error {
	_, err := repository.db.Exec(
		`INSERT into maker.pit_file_stability_fee (header_id, what, data, tx_idx, raw_log)
        VALUES($1, $2, $3, $4, $5)`,
		headerID, model.What, model.Data, model.TransactionIndex, model.Raw,
	)
	return err
}

func (repository PitFileStabilityFeeRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN maker.pit_file_stability_fee on headers.id = header_id
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
