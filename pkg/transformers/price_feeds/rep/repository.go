package rep

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type IRepRepository interface {
	CreateRep(rep price_feeds.PriceUpdate) error
}

type RepRepository struct {
	db *postgres.DB
}

func NewRepRepository(db *postgres.DB) RepRepository {
	return RepRepository{
		db: db,
	}
}

func (repository RepRepository) CreateRep(rep price_feeds.PriceUpdate) error {
	_, err := repository.db.Exec(`INSERT INTO maker.reps (block_number, header_id, usd_value) VALUES ($1, $2, $3::NUMERIC)`, rep.BlockNumber, rep.HeaderID, rep.UsdValue)
	return err
}
