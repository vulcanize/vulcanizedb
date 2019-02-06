package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type MockMappings struct {
	Metadata     shared.StorageValueMetadata
	LookupCalled bool
	LookupErr    error
}

func (mappings *MockMappings) Lookup(key common.Hash) (shared.StorageValueMetadata, error) {
	mappings.LookupCalled = true
	return mappings.Metadata, mappings.LookupErr
}

func (*MockMappings) SetDB(db *postgres.DB) {
	panic("implement me")
}
