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

package dent

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockDentConverter struct {
	converterError        error
	PassedContractAddress string
	PassedContractAbi     string
	LogsToConvert         []types.Log
}

func (c *MockDentConverter) Convert(contractAddress string, contractAbi string, ethLog types.Log) (dent.DentModel, error) {
	c.PassedContractAddress = contractAddress
	c.PassedContractAbi = contractAbi
	c.LogsToConvert = append(c.LogsToConvert, ethLog)
	return test_data.DentModel, c.converterError
}

func (c *MockDentConverter) SetConverterError(err error) {
	c.converterError = err
}
