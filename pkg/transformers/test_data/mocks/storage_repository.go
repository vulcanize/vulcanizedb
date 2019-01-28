package mocks

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type MockStorageRepository struct {
	CreateErr         error
	PassedBlockNumber int
	PassedBlockHash   string
	PassedMetadata    shared.StorageValueMetadata
	PassedValue       interface{}
}

func (repository *MockStorageRepository) Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error {
	repository.PassedBlockNumber = blockNumber
	repository.PassedBlockHash = blockHash
	repository.PassedMetadata = metadata
	repository.PassedValue = value
	return repository.CreateErr
}

func (*MockStorageRepository) SetDB(db *postgres.DB) {
	panic("implement me")
}
