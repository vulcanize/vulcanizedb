package vat_init

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
)

type MockVatInitRepository struct {
	createErr                 error
	missingHeaders            []core.Header
	missingHeadersErr         error
	PassedStartingBlockNumber int64
	PassedEndingBlockNumber   int64
	PassedHeaderID            int64
	PassedModel               vat_init.VatInitModel
}

func (repository *MockVatInitRepository) Create(headerID int64, model vat_init.VatInitModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedModel = model
	return repository.createErr
}

func (repository *MockVatInitRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockVatInitRepository) SetMissingHeadersErr(e error) {
	repository.missingHeadersErr = e
}

func (repository *MockVatInitRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockVatInitRepository) SetCreateError(e error) {
	repository.createErr = e
}
