package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type MockLogNoteConverter struct {
	err          error
	returnModels []interface{}
	PassedLogs   []types.Log
}

func (converter *MockLogNoteConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	converter.PassedLogs = ethLogs
	return converter.returnModels, converter.err
}

func (converter *MockLogNoteConverter) SetConverterError(e error) {
	converter.err = e
}

func (converter *MockLogNoteConverter) SetReturnModels(models []interface{}) {
	converter.returnModels = models
}
