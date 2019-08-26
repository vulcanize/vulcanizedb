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

type MockCheckedLogsRepository struct {
	HasLogBeenCheckedAddresses []string
	HasLogBeenCheckedError     error
	HasLogBeenCheckedReturn    bool
	HasLogBeenCheckedTopicZero string
	MarkLogCheckedAddresses    []string
	MarkLogCheckedError        error
	MarkLogCheckedTopicZero    string
}

func (repository *MockCheckedLogsRepository) HaveLogsBeenChecked(addresses []string, topic0 string) (bool, error) {
	repository.HasLogBeenCheckedAddresses = addresses
	repository.HasLogBeenCheckedTopicZero = topic0
	return repository.HasLogBeenCheckedReturn, repository.HasLogBeenCheckedError
}

func (repository *MockCheckedLogsRepository) MarkLogsChecked(addresses []string, topic0 string) error {
	repository.MarkLogCheckedAddresses = addresses
	repository.MarkLogCheckedTopicZero = topic0
	return repository.MarkLogCheckedError
}
