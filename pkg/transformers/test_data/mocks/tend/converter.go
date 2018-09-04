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

package tend

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockTendConverter struct {
	ConverterContract string
	ConverterAbi      string
	LogsToConvert     []types.Log
	ConverterError    error
}

func (c *MockTendConverter) Convert(contractAddress string, contractAbi string, ethLog types.Log) (tend.TendModel, error) {
	c.ConverterContract = contractAddress
	c.ConverterAbi = contractAbi
	c.LogsToConvert = append(c.LogsToConvert, ethLog)
	return test_data.TendModel, c.ConverterError
}

func (c *MockTendConverter) SetConverterError(err error) {
	c.ConverterError = err
}
