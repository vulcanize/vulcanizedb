// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
