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

package fakes

import (
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type MockCheckedHeadersRepository struct {
	MarkHeaderCheckedHeaderID               int64
	MarkHeaderCheckedReturnError            error
	MarkHeadersUncheckedCalled              bool
	MarkHeadersUncheckedReturnError         error
	MarkHeadersUncheckedStartingBlockNumber int64
	UncheckedHeadersCheckCount              int64
	UncheckedHeadersEndingBlockNumber       int64
	UncheckedHeadersReturnError             error
	UncheckedHeadersReturnHeaders           []core.Header
	UncheckedHeadersStartingBlockNumber     int64
}

func (repository *MockCheckedHeadersRepository) MarkHeadersUnchecked(startingBlockNumber int64) error {
	repository.MarkHeadersUncheckedCalled = true
	repository.MarkHeadersUncheckedStartingBlockNumber = startingBlockNumber
	return repository.MarkHeadersUncheckedReturnError
}

func (repository *MockCheckedHeadersRepository) MarkHeaderChecked(headerID int64) error {
	repository.MarkHeaderCheckedHeaderID = headerID
	return repository.MarkHeaderCheckedReturnError
}

func (repository *MockCheckedHeadersRepository) UncheckedHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error) {
	repository.UncheckedHeadersStartingBlockNumber = startingBlockNumber
	repository.UncheckedHeadersEndingBlockNumber = endingBlockNumber
	repository.UncheckedHeadersCheckCount = checkCount
	return repository.UncheckedHeadersReturnHeaders, repository.UncheckedHeadersReturnError
}
