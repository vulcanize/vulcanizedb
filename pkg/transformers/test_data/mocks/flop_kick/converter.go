// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flop_kick

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockConverter struct {
	PassedContractAddress string
	PassedContractABI     string
	PassedLogs            []types.Log
	PassedEntities        []flop_kick.Entity
	entityConverterError  error
	modelConverterError   error
}

func (c *MockConverter) ToEntities(contractAddress, contractAbi string, ethLogs []types.Log) ([]flop_kick.Entity, error) {
	c.PassedContractAddress = contractAddress
	c.PassedContractABI = contractAbi
	c.PassedLogs = ethLogs

	return []flop_kick.Entity{test_data.FlopKickEntity}, c.entityConverterError
}

func (c *MockConverter) ToModels(entities []flop_kick.Entity) ([]flop_kick.Model, error) {
	c.PassedEntities = entities
	return []flop_kick.Model{test_data.FlopKickModel}, c.modelConverterError
}

func (c *MockConverter) SetToEntityConverterError(err error) {
	c.entityConverterError = err
}

func (c *MockConverter) SetToModelConverterError(err error) {
	c.modelConverterError = err
}
