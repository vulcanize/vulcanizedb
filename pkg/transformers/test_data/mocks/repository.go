package mocks

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockRepository struct {
	createError                      error
	markHeaderCheckedError           error
	MarkHeaderCheckedPassedHeaderIDs []int64
	missingHeaders                   []core.Header
	missingHeadersError              error
	PassedStartingBlockNumber        int64
	PassedEndingBlockNumber          int64
	PassedHeaderID                   int64
	PassedModels                     []interface{}
	SetDbCalled                      bool
}

func (repository *MockRepository) Create(headerID int64, models []interface{}) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	return repository.createError
}

func (repository *MockRepository) MarkHeaderChecked(headerID int64) error {
	repository.MarkHeaderCheckedPassedHeaderIDs = append(repository.MarkHeaderCheckedPassedHeaderIDs, headerID)
	return repository.markHeaderCheckedError
}

func (repository *MockRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}

func (repository *MockRepository) SetMissingHeadersError(e error) {
	repository.missingHeadersError = e
}

func (repository *MockRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.MarkHeaderCheckedPassedHeaderIDs).To(ContainElement(i))
}

func (repository *MockRepository) SetMarkHeaderCheckedError(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockRepository) SetCreateError(e error) {
	repository.createError = e
}
