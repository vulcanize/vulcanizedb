package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/storage"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type MockStorageTransformer struct {
	Address    common.Address
	ExecuteErr error
	PassedRow  shared.StorageDiffRow
}

func (transformer *MockStorageTransformer) Execute(row shared.StorageDiffRow) error {
	transformer.PassedRow = row
	return transformer.ExecuteErr
}

func (transformer *MockStorageTransformer) ContractAddress() common.Address {
	return transformer.Address
}

func (transformer *MockStorageTransformer) FakeTransformerInitializer(db *postgres.DB) storage.Transformer {
	return transformer
}
