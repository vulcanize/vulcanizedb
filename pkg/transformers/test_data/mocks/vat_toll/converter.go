package vat_toll

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
)

type MockVatTollConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatTollConverter) ToModels(ethLogs []types.Log) ([]vat_toll.VatTollModel, error) {
	converter.PassedLogs = ethLogs
	return []vat_toll.VatTollModel{test_data.VatTollModel}, converter.err
}

func (converter *MockVatTollConverter) SetConverterError(e error) {
	converter.err = e
}
