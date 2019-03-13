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

type MockRepository struct {
	createError                      error
	markHeaderCheckedError           error
	MarkHeaderCheckedPassedHeaderIDs []int64
	CreatedHeaderIds                 []int64
	missingHeaders                   []core.Header
	allHeaders                       []core.Header
	missingHeadersError              error
	PassedStartingBlockNumber        int64
	PassedEndingBlockNumber          int64
	PassedHeaderID                   int64
	PassedModels                     []interface{}
	SetDbCalled                      bool
	CreateCalledCounter              int
}

func (repository *MockRepository) Create(headerID int64, models []interface{}) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	repository.CreatedHeaderIds = append(repository.CreatedHeaderIds, headerID)
	repository.CreateCalledCounter++

	return repository.createError
}

func (repository *MockRepository) MarkHeaderChecked(headerID int64) error {
	repository.MarkHeaderCheckedPassedHeaderIDs = append(repository.MarkHeaderCheckedPassedHeaderIDs, headerID)
	return repository.markHeaderCheckedError
}

func (repository *MockRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.missingHeaders, repository.missingHeadersError
}

func (repository *MockRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	repository.PassedStartingBlockNumber = startingBlockNumber
	repository.PassedEndingBlockNumber = endingBlockNumber
	return repository.allHeaders, nil
}

func (repository *MockRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}

func (repository *MockRepository) SetMissingHeadersError(e error) {
	repository.missingHeadersError = e
}

func (repository *MockRepository) SetAllHeaders(headers []core.Header) {
	repository.allHeaders = headers
}

func (repository *MockRepository) SetMissingHeaders(headers []core.Header) {
	repository.missingHeaders = headers
}

func (repository *MockRepository) SetMarkHeaderCheckedError(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockRepository) SetCreateError(e error) {
	repository.createError = e
}

func (repository *MockRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.MarkHeaderCheckedPassedHeaderIDs).To(ContainElement(i))
}

func (repository *MockRepository) AssertMarkHeaderCheckedNotCalled() {
	Expect(len(repository.MarkHeaderCheckedPassedHeaderIDs)).To(Equal(0))
}
