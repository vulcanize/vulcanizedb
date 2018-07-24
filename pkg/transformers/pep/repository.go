package pep

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type IPepRepository interface {
	CreatePep(pep Pep) error
}

type PepRepository struct {
	db *postgres.DB
}

func NewPepRepository(db *postgres.DB) PepRepository {
	return PepRepository{
		db: db,
	}
}

func (repository PepRepository) CreatePep(pep Pep) error {
	_, err := repository.db.Exec(`INSERT INTO maker.peps (block_number, header_id, usd_value) VALUES ($1, $2, $3::NUMERIC)`, pep.BlockNumber, pep.HeaderID, pep.UsdValue)
	if err != nil {
		return err
	}
	return nil
}
