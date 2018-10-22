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

package vow

import (
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockDripFileVowRepository struct {
	createError                     error
	markHeaderCheckedError          error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersError             error
	PassedStartingBlockNumber       int64
	PassedEndingBlockNumber         int64
	PassedHeaderID                  int64
	PassedModels                    []interface{}
	SetDbCalled                     bool
}

func (repository *MockDripFileVowRepository) Create(headerID int64, models []interface{}) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	return repository.createError
}

func (repository *MockDripFileVowRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedError
}

func (repository *MockDripFileVowRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockDripFileVowRepository) SetMarkHeaderCheckedError(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockDripFileVowRepository) SetMissingHeadersError(e error) {
	repository.missingHeadersError = e
}

func (repository *MockDripFileVowRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockDripFileVowRepository) SetCreateError(e error) {
	repository.createError = e
}
func (repository *MockDripFileVowRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(i))
}

func (repository *MockDripFileVowRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}
