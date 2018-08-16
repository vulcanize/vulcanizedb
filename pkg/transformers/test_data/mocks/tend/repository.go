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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/tend"
)

type MockTendRepository struct {
	createError               error
	PassedEndingBlockNumber   int64
	PassedHeaderID            int64
	PassedStartingBlockNumber int64
	PassedTendModel           tend.TendModel
	missingHeaders            []core.Header
	missingHeadersErr         error
}

func (repository *MockTendRepository) Create(headerId int64, tend tend.TendModel) error {
	repository.PassedHeaderID = headerId
	repository.PassedTendModel = tend
	return repository.createError
}

func (repository *MockTendRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockTendRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockTendRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockTendRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}
