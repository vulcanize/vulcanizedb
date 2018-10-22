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
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockPitFileIlkRepository struct {
	createError                     error
	PassedEndingBlockNumber         int64
	PassedModels                    []interface{}
	PassedHeaderID                  int64
	PassedStartingBlockNumber       int64
	markHeaderCheckedError          error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersError             error
	SetDbCalled                     bool
}

func (repository *MockPitFileIlkRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockPitFileIlkRepository) SetMarkHeaderCheckedError(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockPitFileIlkRepository) SetMissingHeadersError(err error) {
	repository.missingHeadersError = err
}

func (repository *MockPitFileIlkRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPitFileIlkRepository) Create(headerID int64, models []interface{}) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	return repository.createError
}

func (repository *MockPitFileIlkRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedError
}

func (repository *MockPitFileIlkRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockPitFileIlkRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(i))
}

func (repository *MockPitFileIlkRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}
