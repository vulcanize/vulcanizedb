// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mocks

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockEventRepository struct {
	createError                      error
	markHeaderCheckedError           error
	MarkHeaderCheckedPassedHeaderIDs []int64
	missingHeaders                   []core.Header
	allHeaders                       []core.Header
	missingHeadersError              error
	PassedStartingBlockNumber        int64
	PassedEndingBlockNumber          int64
	PassedModels                     []interface{}
	SetDbCalled                      bool
	CreateCalledCounter              int
}

func (repository *MockEventRepository) Create(models []interface{}) error {
	repository.PassedModels = models
	repository.CreateCalledCounter++

	return repository.createError
}

func (repository *MockEventRepository) MarkHeaderChecked(headerID int64) error {
	repository.MarkHeaderCheckedPassedHeaderIDs = append(repository.MarkHeaderCheckedPassedHeaderIDs, headerID)
	return repository.markHeaderCheckedError
}

func (repository *MockEventRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockEventRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.allHeaders, nil
}

func (repository *MockEventRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}

func (repository *MockEventRepository) SetMissingHeadersError(e error) {
	repository.missingHeadersError = e
}

func (repository *MockEventRepository) SetAllHeaders(headers []core.Header) {
	repository.allHeaders = headers
}

func (repository *MockEventRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockEventRepository) SetMarkHeaderCheckedError(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockEventRepository) SetCreateError(e error) {
	repository.createError = e
}

func (repository *MockEventRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.MarkHeaderCheckedPassedHeaderIDs).To(ContainElement(i))
}
