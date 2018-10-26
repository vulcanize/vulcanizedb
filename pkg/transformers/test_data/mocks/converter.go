package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type MockConverter struct {
	ToEntitiesError   error
	ToModelsError     error
	ContractAbi       string
	LogsToConvert     []types.Log
	EntitiesToConvert []interface{}
	EntitiesToReturn  []interface{}
	ModelsToReturn    []interface{}
}

func (converter *MockConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	converter.ContractAbi = contractAbi
	converter.LogsToConvert = ethLogs
	return converter.EntitiesToReturn, converter.ToEntitiesError
}

func (converter *MockConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	converter.EntitiesToConvert = entities
	return converter.ModelsToReturn, converter.ToModelsError
}
