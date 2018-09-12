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

package flip

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/cat_file/flip"
)

type MockCatFileFlipRepository struct {
	createErr                       error
	markHeaderCheckedErr            error
	markHeaderCheckedPassedHeaderID int64
	missingHeadersErr               error
	missingHeaders                  []core.Header
	PassedStartingBlockNumber       int64
	PassedEndingBlockNumber         int64
	PassedHeaderID                  int64
	PassedModels                    []flip.CatFileFlipModel
}

func (repository *MockCatFileFlipRepository) Create(headerID int64, models []flip.CatFileFlipModel) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	return repository.createErr
}

func (repository *MockCatFileFlipRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedErr
}

func (repository *MockCatFileFlipRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockCatFileFlipRepository) SetMissingHeadersErr(e error) {
	repository.missingHeadersErr = e
}

func (repository *MockCatFileFlipRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockCatFileFlipRepository) SetCreateError(e error) {
	repository.createErr = e
}

func (repository *MockCatFileFlipRepository) SetMarkHeaderCheckedErr(e error) {
	repository.markHeaderCheckedErr = e
}

func (repository *MockCatFileFlipRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(i))
}
