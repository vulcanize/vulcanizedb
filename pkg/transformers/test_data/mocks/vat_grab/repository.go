package vat_grab

import (
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
)

type MockVatGrabRepository struct {
	createErr                       error
	markHeaderCheckedErr            error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersErr               error
	PassedStartingBlockNumber       int64
	PassedEndingBlockNumber         int64
	PassedHeaderID                  int64
	PassedModels                    []vat_grab.VatGrabModel
}

func (repository *MockVatGrabRepository) Create(headerID int64, models []vat_grab.VatGrabModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	return repository.createErr
}

func (repository *MockVatGrabRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedErr
}

func (repository *MockVatGrabRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockVatGrabRepository) SetMissingHeadersErr(e error) {
	repository.missingHeadersErr = e
}

func (repository *MockVatGrabRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockVatGrabRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(i))
}

func (repository *MockVatGrabRepository) SetMarkHeaderCheckedErr(e error) {
	repository.markHeaderCheckedErr = e
}

func (repository *MockVatGrabRepository) SetCreateError(e error) {
	repository.createErr = e
}
