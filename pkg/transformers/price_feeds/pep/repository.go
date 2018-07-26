package pep

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type IPepRepository interface {
	CreatePep(pep price_feeds.PriceUpdate) error
}

type PepRepository struct {
	db *postgres.DB
}

func NewPepRepository(db *postgres.DB) PepRepository {
	return PepRepository{
		db: db,
	}
}

func (repository PepRepository) CreatePep(pep price_feeds.PriceUpdate) error {
	_, err := repository.db.Exec(`INSERT INTO maker.peps (block_number, header_id, usd_value) VALUES ($1, $2, $3::NUMERIC)`, pep.BlockNumber, pep.HeaderID, pep.UsdValue)
	return err
}
