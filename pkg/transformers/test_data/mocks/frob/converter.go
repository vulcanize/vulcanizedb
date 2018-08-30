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
	PassedContractAddress string
	PassedContractABI     string
	PassedLog             types.Log
	PassedEntity          frob.FrobEntity
	toEntityError         error
	toModelError          error
}

func (converter *MockFrobConverter) SetToEntityError(err error) {
	converter.toEntityError = err
}

func (converter *MockFrobConverter) SetToModelError(err error) {
	converter.toModelError = err
}

func (converter *MockFrobConverter) ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (frob.FrobEntity, error) {
	converter.PassedContractAddress = contractAddress
	converter.PassedContractABI = contractAbi
	converter.PassedLog = ethLog
	return test_data.FrobEntity, converter.toEntityError
}

func (converter *MockFrobConverter) ToModel(frobEntity frob.FrobEntity) (frob.FrobModel, error) {
	converter.PassedEntity = frobEntity
	return test_data.FrobModel, converter.toModelError
}
