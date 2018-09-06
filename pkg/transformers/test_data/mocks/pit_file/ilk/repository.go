package ilk

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
)

type MockPitFileIlkRepository struct {
	createError               error
	PassedEndingBlockNumber   int64
	PassedModel               ilk.PitFileIlkModel
	PassedHeaderID            int64
	PassedStartingBlockNumber int64
	missingHeaders            []core.Header
	missingHeadersErr         error
}

func (repository *MockPitFileIlkRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockPitFileIlkRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockPitFileIlkRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPitFileIlkRepository) Create(headerID int64, model ilk.PitFileIlkModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedModel = model
	return repository.createError
}

func (repository *MockPitFileIlkRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}
