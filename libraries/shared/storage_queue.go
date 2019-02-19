package shared

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type IStorageQueue interface {
	Add(row shared.StorageDiffRow) error
}

type StorageQueue struct {
	db *postgres.DB
}

func NewStorageQueue(db *postgres.DB) StorageQueue {
	return StorageQueue{db: db}
}

func (queue StorageQueue) Add(row shared.StorageDiffRow) error {
	_, err := queue.db.Exec(`INSERT INTO public.queued_storage (contract,
		block_hash, block_height, storage_key, storage_value) VALUES
		($1, $2, $3, $4, $5)`, row.Contract.Bytes(), row.BlockHash.Bytes(),
		row.BlockHeight, row.StorageKey.Bytes(), row.StorageValue.Bytes())
	return err
}
