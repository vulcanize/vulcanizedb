package vat_init

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
)

type MockVatInitConverter struct {
	converterErr error
	PassedLog    types.Log
}

func (converter *MockVatInitConverter) ToModel(ethLog types.Log) (vat_init.VatInitModel, error) {
	converter.PassedLog = ethLog
	return test_data.VatInitModel, converter.converterErr
}

func (converter *MockVatInitConverter) SetConverterError(e error) {
	converter.converterErr = e
}
