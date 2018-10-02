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

package frob

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockFrobConverter struct {
	PassedContractABI string
	PassedLogs        []types.Log
	PassedEntities    []frob.FrobEntity
	toEntityError     error
	toModelError      error
}

func (converter *MockFrobConverter) SetToEntitiesError(err error) {
	converter.toEntityError = err
}

func (converter *MockFrobConverter) SetToModelsError(err error) {
	converter.toModelError = err
}

func (converter *MockFrobConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]frob.FrobEntity, error) {
	converter.PassedContractABI = contractAbi
	converter.PassedLogs = ethLogs
	return []frob.FrobEntity{test_data.FrobEntity}, converter.toEntityError
}

func (converter *MockFrobConverter) ToModels(frobEntities []frob.FrobEntity) ([]frob.FrobModel, error) {
	converter.PassedEntities = frobEntities
	return []frob.FrobModel{test_data.FrobModel}, converter.toModelError
}
