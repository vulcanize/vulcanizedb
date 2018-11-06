package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type MockConverter struct {
	ToEntitiesError         error
	PassedContractAddresses []string
	ToModelsError           error
	entityConverterError    error
	modelConverterError     error
	ContractAbi             string
	LogsToConvert           []types.Log
	EntitiesToConvert       []interface{}
	EntitiesToReturn        []interface{}
	ModelsToReturn          []interface{}
	ToEntitiesCalledCounter int
	ToModelsCalledCounter   int
}

func (converter *MockConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	for _, log := range ethLogs {
		converter.PassedContractAddresses = append(converter.PassedContractAddresses, log.Address.Hex())
	}
	converter.ContractAbi = contractAbi
	converter.LogsToConvert = ethLogs
	return converter.EntitiesToReturn, converter.ToEntitiesError
}

func (converter *MockConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	converter.EntitiesToConvert = entities
	return converter.ModelsToReturn, converter.ToModelsError
}

func (converter *MockConverter) SetToEntityConverterError(err error) {
	converter.entityConverterError = err
}

func (c *MockConverter) SetToModelConverterError(err error) {
	c.modelConverterError = err
}
