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

package datastore

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type AddressRepository interface {
	GetOrCreateAddress(address string) (int, error)
}

type CheckedHeadersRepository interface {
	MarkHeaderChecked(headerID int64) error
	MarkHeadersUnchecked(startingBlockNumber int64) error
	MarkSingleHeaderUnchecked(blockNumber int64) error
	UncheckedHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error)
}

type CheckedLogsRepository interface {
	AlreadyWatchingLog(addresses []string, topic0 string) (bool, error)
	MarkLogWatched(addresses []string, topic0 string) error
}

type HeaderRepository interface {
	CreateOrUpdateHeader(header core.Header) (int64, error)
	CreateTransactions(headerID int64, transactions []core.TransactionModel) error
	GetHeader(blockNumber int64) (core.Header, error)
	MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64) ([]int64, error)
}

type EventLogRepository interface {
	GetUntransformedEventLogs() ([]core.EventLog, error)
	CreateEventLogs(headerID int64, logs []types.Log) error
}
