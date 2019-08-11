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

import "github.com/vulcanize/vulcanizedb/pkg/core"

type MockCheckedHeadersRepository struct {
	CheckCount                   int64
	StartingBlockNumber          int64
	EndingBlockNumber            int64
	HeaderID                     int64
	ReturnHeaders                []core.Header
	MarkHeaderCheckedReturnError error
	MissingHeadersReturnError    error
}

func (repository *MockCheckedHeadersRepository) MarkHeaderChecked(headerID int64) error {
	repository.HeaderID = headerID
	return repository.MarkHeaderCheckedReturnError
}

func (repository *MockCheckedHeadersRepository) MissingHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error) {
	repository.StartingBlockNumber = startingBlockNumber
	repository.EndingBlockNumber = endingBlockNumber
	repository.CheckCount = checkCount
	return repository.ReturnHeaders, repository.MissingHeadersReturnError
}
