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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockHeaderSyncLogRepository struct {
	CreateError    error
	GetCalled      bool
	GetError       error
	PassedHeaderID int64
	PassedLogs     []types.Log
	ReturnLogs     []core.HeaderSyncLog
}

func (repository *MockHeaderSyncLogRepository) GetUntransformedHeaderSyncLogs() ([]core.HeaderSyncLog, error) {
	repository.GetCalled = true
	return repository.ReturnLogs, repository.GetError
}

func (repository *MockHeaderSyncLogRepository) CreateHeaderSyncLogs(headerID int64, logs []types.Log) error {
	repository.PassedHeaderID = headerID
	repository.PassedLogs = logs
	return repository.CreateError
}
