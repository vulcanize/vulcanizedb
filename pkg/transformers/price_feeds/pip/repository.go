package pip

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type IPipRepository interface {
	CreatePip(pip price_feeds.PriceUpdate) error
}

type PipRepository struct {
	db *postgres.DB
}

func NewPipRepository(db *postgres.DB) PipRepository {
	return PipRepository{
		db: db,
	}
}

func (repository PipRepository) CreatePip(pip price_feeds.PriceUpdate) error {
	_, err := repository.db.Exec(`INSERT INTO maker.pips (block_number, header_id, usd_value) VALUES ($1, $2, $3::NUMERIC)`, pip.BlockNumber, pip.HeaderID, pip.UsdValue)
	return err
}
