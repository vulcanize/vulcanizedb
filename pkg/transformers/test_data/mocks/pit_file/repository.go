package pit_file

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
)

type MockPitFileRepository struct {
	createError               error
	PassedEndingBlockNumber   int64
	PassedModel               pit_file.PitFileModel
	PassedHeaderID            int64
	PassedStartingBlockNumber int64
	PassedTransactionIndex    uint
	missingHeaders            []core.Header
	missingHeadersErr         error
}

func (repository *MockPitFileRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockPitFileRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockPitFileRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPitFileRepository) Create(headerID int64, transactionIndex uint, model pit_file.PitFileModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedTransactionIndex = transactionIndex
	repository.PassedModel = model
	return repository.createError
}

func (repository *MockPitFileRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}
