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

package flop_kick

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
)

type MockRepository struct {
	PassedStartingBlockNumber int64
	PassedEndingBlockNumber   int64
	CreatedHeaderIds          []int64
	CreatedModels             []flop_kick.Model
	CheckedHeaderIds          []int64
	missingHeaders            []core.Header
	missingHeadersError       error
	createError               error
	checkedHeaderError        error
}

func (r *MockRepository) Create(headerId int64, flopKicks []flop_kick.Model) error {
	r.CreatedHeaderIds = append(r.CreatedHeaderIds, headerId)
	r.CreatedModels = flopKicks

	return r.createError
}

func (r *MockRepository) MarkHeaderChecked(headerId int64) error {
	r.CheckedHeaderIds = append(r.CheckedHeaderIds, headerId)

	return r.checkedHeaderError
}

func (r *MockRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	r.PassedStartingBlockNumber = startingBlockNumber
	r.PassedEndingBlockNumber = endingBlockNumber

	return r.missingHeaders, r.missingHeadersError
}

func (r *MockRepository) SetMissingHeaders(headers []core.Header) {
	r.missingHeaders = headers
}

func (r *MockRepository) SetMissingHeadersError(err error) {
	r.missingHeadersError = err
}

func (r *MockRepository) SetCreateError(err error) {
	r.createError = err
}

func (r *MockRepository) SetCheckedHeaderError(err error) {
	r.checkedHeaderError = err
}
