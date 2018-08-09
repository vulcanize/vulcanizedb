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

package frob

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/frob"
)

type MockFrobRepository struct {
	createError               error
	PassedEndingBlockNumber   int64
	PassedFrobModel           frob.FrobModel
	PassedHeaderID            int64
	PassedStartingBlockNumber int64
	PassedTransactionIndex    uint
	missingHeaders            []core.Header
	missingHeadersErr         error
}

func (repository *MockFrobRepository) SetCreateError(err error) {
	repository.createError = err
}

func (repository *MockFrobRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockFrobRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockFrobRepository) Create(headerID int64, transactionIndex uint, model frob.FrobModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedTransactionIndex = transactionIndex
	repository.PassedFrobModel = model
	return repository.createError
}

func (repository *MockFrobRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}
