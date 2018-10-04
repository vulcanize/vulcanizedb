package vat_tune

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
)

type MockVatTuneConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatTuneConverter) ToModels(ethLogs []types.Log) ([]vat_tune.VatTuneModel, error) {
	converter.PassedLogs = ethLogs
	return []vat_tune.VatTuneModel{test_data.VatTuneModel}, converter.err
}

func (converter *MockVatTuneConverter) SetConverterError(e error) {
	converter.err = e
}
