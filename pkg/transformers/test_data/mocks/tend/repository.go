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
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockTendRepository struct {
	createError                     error
	PassedEndingBlockNumber         int64
	PassedHeaderID                  int64
	PassedStartingBlockNumber       int64
	PassedTendModel                 interface{}
	markHeaderCheckedError          error
	markHeaderCheckedPassedHeaderId int64
	missingHeaders                  []core.Header
	missingHeadersError             error
	SetDbCalled                     bool
}

func (repository *MockTendRepository) Create(headerId int64, tend []interface{}) error {
	repository.PassedHeaderID = headerId
	repository.PassedTendModel = tend[0]
	return repository.createError
}

func (repository *MockTendRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockTendRepository) SetMarkHeaderCheckedErr(err error) {
	repository.markHeaderCheckedError = err
}

func (repository *MockTendRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersError = err
}

func (repository *MockTendRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockTendRepository) MarkHeaderChecked(headerId int64) error {
	repository.markHeaderCheckedPassedHeaderId = headerId
	return repository.markHeaderCheckedError
}

func (repository *MockTendRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockTendRepository) AssertMarkHeaderCheckedCalledWith(headerId int64) {
	Expect(repository.markHeaderCheckedPassedHeaderId).To(Equal(headerId))
}

func (repository *MockTendRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}
