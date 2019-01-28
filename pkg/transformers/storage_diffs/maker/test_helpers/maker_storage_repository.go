package test_helpers

import "github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

type MockMakerStorageRepository struct {
	GetIlksCalled bool
	ilks          []string
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.ilks, nil
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}

func (repository *MockMakerStorageRepository) SetIlks(ilks []string) {
	repository.ilks = ilks
}
