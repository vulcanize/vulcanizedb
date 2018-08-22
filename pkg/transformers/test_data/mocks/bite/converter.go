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

package bite

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/vulcanize/vulcanizedb/pkg/transformers/bite"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockBiteConverter struct {
	ConverterContract string
	ConverterAbi      string
	LogsToConvert     []types.Log
	EntitiesToConvert []BiteEntity
	ConverterError    error
}

func (mbc *MockBiteConverter) ToEntity(contractAddress string, contractAbi string, ethLog types.Log) (BiteEntity, error) {
	mbc.ConverterContract = contractAddress
	mbc.ConverterAbi = contractAbi
	mbc.LogsToConvert = append(mbc.LogsToConvert, ethLog)
	return test_data.BiteEntity, mbc.ConverterError
}

func (mbc *MockBiteConverter) ToModel(entity BiteEntity) (BiteModel, error) {
	mbc.EntitiesToConvert = append(mbc.EntitiesToConvert, entity)
	return test_data.BiteModel, mbc.ConverterError
}
func (mbc *MockBiteConverter) SetConverterError(err error) {
	mbc.ConverterError = err
}
