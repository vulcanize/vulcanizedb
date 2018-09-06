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

package vat_init

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
)

type MockVatInitConverter struct {
	converterErr error
	PassedLog    types.Log
}

func (converter *MockVatInitConverter) ToModel(ethLog types.Log) (vat_init.VatInitModel, error) {
	converter.PassedLog = ethLog
	return test_data.VatInitModel, converter.converterErr
}

func (converter *MockVatInitConverter) SetConverterError(e error) {
	converter.converterErr = e
}
