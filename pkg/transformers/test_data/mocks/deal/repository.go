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

package deal

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
)

type MockDealRepository struct {
	createError                     error
	PassedEndingBlockNumber         int64
	PassedHeaderIDs                 []int64
	PassedStartingBlockNumber       int64
	PassedDealModels                []deal.DealModel
	markHeaderCheckedErr            error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersErr               error
}

func (repository *MockDealRepository) Create(headerId int64, deals []deal.DealModel) error {
	repository.PassedHeaderIDs = append(repository.PassedHeaderIDs, headerId)
	repository.PassedDealModels = append(repository.PassedDealModels, deals...)
	return repository.createError
}

func (repository *MockDealRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockDealRepository) SetMarkHeaderCheckedErr(err error) {
	repository.markHeaderCheckedErr = err
}

func (repository *MockDealRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockDealRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockDealRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedErr
}

func (repository *MockDealRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockDealRepository) AssertMarkHeaderCheckedCalledWith(headerID int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(headerID))
}
