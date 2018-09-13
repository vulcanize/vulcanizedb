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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/deal"
)

type MockDealRepository struct {
	createError               error
	PassedEndingBlockNumber   int64
	PassedHeaderIDs           []int64
	PassedStartingBlockNumber int64
	PassedDealModels          []deal.DealModel
	missingHeaders            []core.Header
	missingHeadersErr         error
}

func (repository *MockDealRepository) Create(headerId int64, deal deal.DealModel) error {
	repository.PassedHeaderIDs = append(repository.PassedHeaderIDs, headerId)
	repository.PassedDealModels = append(repository.PassedDealModels, deal)
	return repository.createError
}

func (repository *MockDealRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockDealRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockDealRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockDealRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}
