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

package ilk

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
)

type MockPitFileIlkRepository struct {
	createError               error
	PassedEndingBlockNumber   int64
	PassedModel               ilk.PitFileIlkModel
	PassedHeaderID            int64
	PassedStartingBlockNumber int64
	missingHeaders            []core.Header
	missingHeadersErr         error
}

func (repository *MockPitFileIlkRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockPitFileIlkRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockPitFileIlkRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPitFileIlkRepository) Create(headerID int64, model ilk.PitFileIlkModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedModel = model
	return repository.createError
}

func (repository *MockPitFileIlkRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}
