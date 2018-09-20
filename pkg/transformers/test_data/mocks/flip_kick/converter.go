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

package flip_kick

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockFlipKickConverter struct {
	ConverterContract string
	ConverterAbi      string
	LogsToConvert     []types.Log
	EntitiesToConvert []flip_kick.FlipKickEntity
	ConverterError    error
}

func (mfkc *MockFlipKickConverter) ToEntities(contractAddress string, contractAbi string, ethLogs []types.Log) ([]flip_kick.FlipKickEntity, error) {
	mfkc.ConverterContract = contractAddress
	mfkc.ConverterAbi = contractAbi
	mfkc.LogsToConvert = append(mfkc.LogsToConvert, ethLogs...)
	return []flip_kick.FlipKickEntity{test_data.FlipKickEntity}, mfkc.ConverterError
}

func (mfkc *MockFlipKickConverter) ToModels(flipKickEntities []flip_kick.FlipKickEntity) ([]flip_kick.FlipKickModel, error) {
	mfkc.EntitiesToConvert = append(mfkc.EntitiesToConvert, flipKickEntities...)
	return []flip_kick.FlipKickModel{test_data.FlipKickModel}, nil
}

func (mfkc *MockFlipKickConverter) SetConverterError(err error) {
	mfkc.ConverterError = err
}
