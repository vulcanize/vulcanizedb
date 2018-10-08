package vat_grab

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
)

type MockVatGrabConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatGrabConverter) ToModels(ethLogs []types.Log) ([]vat_grab.VatGrabModel, error) {
	converter.PassedLogs = ethLogs
	return []vat_grab.VatGrabModel{test_data.VatGrabModel}, converter.err
}

func (converter *MockVatGrabConverter) SetConverterError(e error) {
	converter.err = e
}
