package test_helpers

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"math/big"
)

type MockMakerStorageRepository struct {
	DaiKeys          []string
	GemKeys          []maker.Urn
	GetDaiKeysCalled bool
	GetDaiKeysError  error
	GetGemKeysCalled bool
	GetGemKeysError  error
	GetIlksCalled    bool
	GetIlksError     error
	GetMaxFlipCalled bool
	GetMaxFlipError  error
	GetSinKeysCalled bool
	GetSinKeysError  error
	GetUrnsCalled    bool
	GetUrnsError     error
	Ilks             []string
	MaxFlip          *big.Int
	SinKeys          []string
	Urns             []maker.Urn
}

func (repository *MockMakerStorageRepository) GetDaiKeys() ([]string, error) {
	repository.GetDaiKeysCalled = true
	return repository.DaiKeys, repository.GetDaiKeysError
}

func (repository *MockMakerStorageRepository) GetGemKeys() ([]maker.Urn, error) {
	repository.GetGemKeysCalled = true
	return repository.GemKeys, repository.GetGemKeysError
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.Ilks, repository.GetIlksError
}

func (repository *MockMakerStorageRepository) GetMaxFlip() (*big.Int, error) {
	repository.GetMaxFlipCalled = true
	return repository.MaxFlip, repository.GetMaxFlipError
}

func (repository *MockMakerStorageRepository) GetSinKeys() ([]string, error) {
	repository.GetSinKeysCalled = true
	return repository.SinKeys, repository.GetSinKeysError
}

func (repository *MockMakerStorageRepository) GetUrns() ([]maker.Urn, error) {
	repository.GetUrnsCalled = true
	return repository.Urns, repository.GetUrnsError
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}
