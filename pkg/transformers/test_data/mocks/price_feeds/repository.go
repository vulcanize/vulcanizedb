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

package price_feeds

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
)

type MockPriceFeedRepository struct {
	createErr                       error
	createPassedHeaderID            int64
	markHeaderCheckedErr            error
	markHeaderCheckedPassedHeaderID int64
	missingHeaders                  []core.Header
	missingHeadersErr               error
	passedEndingBlockNumber         int64
	passedModel                     price_feeds.PriceFeedModel
	passedStartingBlockNumber       int64
}

func (repository *MockPriceFeedRepository) SetCreateErr(err error) {
	repository.createErr = err
}

func (repository *MockPriceFeedRepository) SetMarkHeaderCheckedErr(err error) {
	repository.markHeaderCheckedErr = err
}

func (repository *MockPriceFeedRepository) SetMissingHeadersErr(err error) {
	repository.missingHeadersErr = err
}

func (repository *MockPriceFeedRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockPriceFeedRepository) Create(headerID int64, model price_feeds.PriceFeedModel) error {
	repository.createPassedHeaderID = headerID
	repository.passedModel = model
	return repository.createErr
}

func (repository *MockPriceFeedRepository) MarkHeaderChecked(headerID int64) error {
	repository.markHeaderCheckedPassedHeaderID = headerID
	return repository.markHeaderCheckedErr
}

func (repository *MockPriceFeedRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.passedStartingBlockNumber = startingBlockNumber
	repository.passedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersErr
}

func (repository *MockPriceFeedRepository) AssertCreateCalledWith(headerID int64, model price_feeds.PriceFeedModel) {
	Expect(repository.createPassedHeaderID).To(Equal(headerID))
	Expect(repository.passedModel).To(Equal(model))
}

func (repository *MockPriceFeedRepository) AssertMarkHeaderCheckedCalledWith(headerID int64) {
	Expect(repository.markHeaderCheckedPassedHeaderID).To(Equal(headerID))
}

func (repository *MockPriceFeedRepository) AssertMissingHeadersCalledwith(startingBlockNumber, endingBlockNumber int64) {
	Expect(repository.passedStartingBlockNumber).To(Equal(startingBlockNumber))
	Expect(repository.passedEndingBlockNumber).To(Equal(endingBlockNumber))
}
