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

package dent

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/dent"
)

type MockDentRepository struct {
	PassedStartingBlockNumber int64
	PassedEndingBlockNumber   int64
	PassedDentModels          []dent.DentModel
	PassedHeaderIds           []int64
	missingHeaders            []core.Header
	missingHeadersError       error
	createError               error
}

func (r *MockDentRepository) Create(headerId int64, model dent.DentModel) error {
	r.PassedHeaderIds = append(r.PassedHeaderIds, headerId)
	r.PassedDentModels = append(r.PassedDentModels, model)

	return r.createError
}

func (r *MockDentRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	r.PassedStartingBlockNumber = startingBlockNumber
	r.PassedEndingBlockNumber = endingBlockNumber

	return r.missingHeaders, r.missingHeadersError
}

func (r *MockDentRepository) SetMissingHeadersError(err error) {
	r.missingHeadersError = err
}

func (r *MockDentRepository) SetMissingHeaders(headers []core.Header) {
	r.missingHeaders = headers
}

func (r *MockDentRepository) SetCreateError(err error) {
	r.createError = err
}
