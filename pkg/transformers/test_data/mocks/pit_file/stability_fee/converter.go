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

package stability_fee

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockPitFileStabilityFeeConverter struct {
	converterErr error
	PassedLog    types.Log
}

func (converter *MockPitFileStabilityFeeConverter) ToModel(ethLog types.Log) (stability_fee.PitFileStabilityFeeModel, error) {
	converter.PassedLog = ethLog
	return test_data.PitFileStabilityFeeModel, converter.converterErr
}
func (converter *MockPitFileStabilityFeeConverter) SetConverterError(e error) {
	converter.converterErr = e
}
