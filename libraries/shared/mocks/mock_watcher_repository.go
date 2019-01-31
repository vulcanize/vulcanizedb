// VulcanizeDB
// Copyright © 2018 Vulcanize

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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockWatcherRepository struct {
	ReturnCheckedColumnNames    []string
	GetCheckedColumnNamesError  error
	GetCheckedColumnNamesCalled bool

	ReturnNotCheckedSQL       string
	CreateNotCheckedSQLCalled bool

	ReturnMissingHeaders []core.Header
	MissingHeadersError  error
	MissingHeadersCalled bool
}

func (repository *MockWatcherRepository) GetCheckedColumnNames(db *postgres.DB) ([]string, error) {
	repository.GetCheckedColumnNamesCalled = true
	if repository.GetCheckedColumnNamesError != nil {
		return []string{}, repository.GetCheckedColumnNamesError
	}

	return repository.ReturnCheckedColumnNames, nil
}

func (repository *MockWatcherRepository) SetCheckedColumnNames(checkedColumnNames []string) {
	repository.ReturnCheckedColumnNames = checkedColumnNames
}

func (repository *MockWatcherRepository) CreateNotCheckedSQL(boolColumns []string) string {
	repository.CreateNotCheckedSQLCalled = true
	return repository.ReturnNotCheckedSQL
}

func (repository *MockWatcherRepository) SetNotCheckedSQL(notCheckedSql string) {
	repository.ReturnNotCheckedSQL = notCheckedSql
}

func (repository *MockWatcherRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, db *postgres.DB, notCheckedSQL string) ([]core.Header, error) {
	if repository.MissingHeadersError != nil {
		return []core.Header{}, repository.MissingHeadersError
	}
	repository.MissingHeadersCalled = true
	return repository.ReturnMissingHeaders, nil
}

func (repository *MockWatcherRepository) SetMissingHeaders(headers []core.Header) {
	repository.ReturnMissingHeaders = headers
}
