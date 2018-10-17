package vat_flux

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
)

type MockVatFlux struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatFlux) ToModels(ethLogs []types.Log) ([]vat_flux.VatFluxModel, error) {
	converter.PassedLogs = ethLogs
	return []vat_flux.VatFluxModel{test_data.VatFluxModel}, converter.err
}

func (converter *MockVatFlux) SetConverterError(e error) {
	converter.err = e
}
