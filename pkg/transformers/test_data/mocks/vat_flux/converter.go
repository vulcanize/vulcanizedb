package vat_flux

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockVatFluxConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatFluxConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	converter.PassedLogs = ethLogs
	return []interface{}{test_data.VatFluxModel}, converter.err
}

func (converter *MockVatFluxConverter) SetConverterError(e error) {
	converter.err = e
}
