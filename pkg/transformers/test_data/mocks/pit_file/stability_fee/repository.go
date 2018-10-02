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

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
)

type MockPitFileStabilityFeeRepository struct {
	createErr                       error
	markHeaderCheckedErr            error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersErr               error
	PassedStartingBlockNumber       int64
	PassedEndingBlockNumber         int64
	PassedHeaderID                  int64
	PassedModels                    []stability_fee.PitFileStabilityFeeModel
}

func (repository *MockPitFileStabilityFeeRepository) Create(headerID int64, models []stability_fee.PitFileStabilityFeeModel) error {
	repository.PassedModels = models
	repository.PassedHeaderID = headerID
	return repository.createErr
}

func (repository *MockPitFileStabilityFeeRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedErr
}

func (repository *MockPitFileStabilityFeeRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockPitFileStabilityFeeRepository) SetMarkHeaderCheckedErr(e error) {
	repository.markHeaderCheckedErr = e
}

func (repository *MockPitFileStabilityFeeRepository) SetMissingHeadersErr(e error) {
	repository.missingHeadersErr = e
}

func (repository *MockPitFileStabilityFeeRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPitFileStabilityFeeRepository) SetCreateError(e error) {
	repository.createErr = e
}

func (repository *MockPitFileStabilityFeeRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(i))
}
