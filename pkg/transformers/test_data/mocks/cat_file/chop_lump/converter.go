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

package chop_lump

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/chop_lump"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockCatFileChopLumpConverter struct {
	Err        error
	PassedLogs []types.Log
}

func (converter *MockCatFileChopLumpConverter) ToModels(ethLogs []types.Log) ([]chop_lump.CatFileChopLumpModel, error) {
	converter.PassedLogs = ethLogs
	return []chop_lump.CatFileChopLumpModel{test_data.CatFileChopLumpModel}, converter.Err
}

func (converter *MockCatFileChopLumpConverter) SetConverterError(e error) {
	converter.Err = e
}
