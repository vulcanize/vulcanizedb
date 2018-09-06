package stability_fee

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
)

type MockPitFileStabilityFeeRepository struct {
	createErr                 error
	missingHeaders            []core.Header
	missingHeadersErr         error
	PassedStartingBlockNumber int64
	PassedEndingBlockNumber   int64
	PassedHeaderID            int64
	PassedModel               stability_fee.PitFileStabilityFeeModel
}

func (repository *MockPitFileStabilityFeeRepository) Create(headerID int64, model stability_fee.PitFileStabilityFeeModel) error {
	repository.PassedModel = model
	repository.PassedHeaderID = headerID
	return repository.createErr
}

func (repository *MockPitFileStabilityFeeRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockPitFileStabilityFeeRepository) SetMissingHeadersErr(e error) {
	repository.missingHeadersErr = e
}
func (repository *MockPitFileStabilityFeeRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}
func (repository *MockPitFileStabilityFeeRepository) SetCreateError(e error) {
	repository.createErr = e
}
