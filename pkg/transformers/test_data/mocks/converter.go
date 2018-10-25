package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type MockConverter struct {
	err          error
	returnModels []interface{}
	PassedLogs   []types.Log
}

func (converter *MockConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	converter.PassedLogs = ethLogs
	return converter.returnModels, converter.err
}

func (converter *MockConverter) SetConverterError(e error) {
	converter.err = e
}

func (converter *MockConverter) SetReturnModels(models []interface{}) {
	converter.returnModels = models
}
