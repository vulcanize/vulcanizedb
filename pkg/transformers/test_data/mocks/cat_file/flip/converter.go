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

package flip

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockCatFileFlipConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockCatFileFlipConverter) ToModels(ethLogs []types.Log) ([]flip.CatFileFlipModel, error) {
	converter.PassedLogs = ethLogs
	return []flip.CatFileFlipModel{test_data.CatFileFlipModel}, converter.err
}

func (converter *MockCatFileFlipConverter) SetConverterError(e error) {
	converter.err = e
}
