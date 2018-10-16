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
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockPitFileStabilityFeeRepository struct {
	createError                     error
	markHeaderCheckedError          error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersError             error
	PassedStartingBlockNumber       int64
	PassedEndingBlockNumber         int64
	PassedHeaderID                  int64
	PassedModels                    []interface{}
}

func (repository *MockPitFileStabilityFeeRepository) Create(headerID int64, models []interface{}) error {
	repository.PassedModels = models
	repository.PassedHeaderID = headerID
	return repository.createError
}

func (repository *MockPitFileStabilityFeeRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedError
}

func (repository *MockPitFileStabilityFeeRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockPitFileStabilityFeeRepository) SetMarkHeaderCheckedErr(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockPitFileStabilityFeeRepository) SetMissingHeadersErr(e error) {
	repository.missingHeadersError = e
}

func (repository *MockPitFileStabilityFeeRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPitFileStabilityFeeRepository) SetCreateError(e error) {
	repository.createError = e
}

func (repository *MockPitFileStabilityFeeRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(i))
}

func (repository *MockPitFileStabilityFeeRepository) SetDB(db *postgres.DB) {}
