package vat_toll

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockVatTollConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockVatTollConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	converter.PassedLogs = ethLogs
	return []interface{}{test_data.VatTollModel}, converter.err
}

func (converter *MockVatTollConverter) SetConverterError(e error) {
	converter.err = e
}
