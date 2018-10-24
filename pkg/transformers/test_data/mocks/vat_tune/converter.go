package vat_tune

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockVatTuneConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatTuneConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	converter.PassedLogs = ethLogs
	return []interface{}{test_data.VatTuneModel}, converter.err
}

func (converter *MockVatTuneConverter) SetConverterError(e error) {
	converter.err = e
}
