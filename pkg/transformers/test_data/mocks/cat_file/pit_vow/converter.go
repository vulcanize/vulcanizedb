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

package pit_vow

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/pit_vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockCatFilePitVowConverter struct {
	err        error
	PassedLogs []types.Log
}

func (converter *MockCatFilePitVowConverter) ToModels(ethLogs []types.Log) ([]pit_vow.CatFilePitVowModel, error) {
	converter.PassedLogs = ethLogs
	return []pit_vow.CatFilePitVowModel{test_data.CatFilePitVowModel}, converter.err
}

func (converter *MockCatFilePitVowConverter) SetConverterError(e error) {
	converter.err = e
}
