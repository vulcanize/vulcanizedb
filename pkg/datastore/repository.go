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
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type AddressRepository interface {
	GetOrCreateAddress(address string) (int, error)
}

type BlockRepository interface {
	CreateOrUpdateBlock(block core.Block) (int64, error)
	GetBlock(blockNumber int64) (core.Block, error)
	MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64
	SetBlocksStatus(chainHead int64) error
}

type CheckedHeadersRepository interface {
	MarkHeaderChecked(headerID int64) error
	MarkHeadersUnchecked(startingBlockNumber int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber, checkCount int64) ([]core.Header, error)
}

type CheckedLogsRepository interface {
	HaveLogsBeenChecked(addresses []string, topic0 string) (bool, error)
	MarkLogsChecked(addresses []string, topic0 string) error
}

type ContractRepository interface {
	CreateContract(contract core.Contract) error
	GetContract(contractHash string) (core.Contract, error)
	ContractExists(contractHash string) (bool, error)
}

type FilterRepository interface {
	CreateFilter(filter filters.LogFilter) error
	GetFilter(name string) (filters.LogFilter, error)
}

type FullSyncLogRepository interface {
	CreateLogs(logs []core.FullSyncLog, receiptId int64) error
	GetLogs(address string, blockNumber int64) ([]core.FullSyncLog, error)
}

type HeaderRepository interface {
	CreateOrUpdateHeader(header core.Header) (int64, error)
	CreateTransactions(headerID int64, transactions []core.TransactionModel) error
	GetHeader(blockNumber int64) (core.Header, error)
	MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) ([]int64, error)
}

type HeaderSyncLogRepository interface {
	GetUntransformedHeaderSyncLogs() ([]core.HeaderSyncLog, error)
	CreateHeaderSyncLogs(headerID int64, logs []types.Log) error
}

type FullSyncReceiptRepository interface {
	CreateReceiptsAndLogs(blockId int64, receipts []core.Receipt) error
	CreateFullSyncReceiptInTx(blockId int64, receipt core.Receipt, tx *sqlx.Tx) (int64, error)
	GetFullSyncReceipt(txHash string) (core.Receipt, error)
}

type HeaderSyncReceiptRepository interface {
	CreateFullSyncReceiptInTx(blockId int64, receipt core.Receipt, tx *sqlx.Tx) (int64, error)
}

type WatchedEventRepository interface {
	GetWatchedEvents(name string) ([]*core.WatchedEvent, error)
}
