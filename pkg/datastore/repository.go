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

package datastore

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var ErrBlockDoesNotExist = func(blockNumber int64) error {
	return fmt.Errorf("Block number %d does not exist", blockNumber)
}

type BlockRepository interface {
	CreateOrUpdateBlock(block core.Block) (int64, error)
	GetBlock(blockNumber int64) (core.Block, error)
	MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64
	SetBlocksStatus(chainHead int64)
}

var ErrContractDoesNotExist = func(contractHash string) error {
	return fmt.Errorf("Contract %v does not exist", contractHash)
}

type ContractRepository interface {
	CreateContract(contract core.Contract) error
	GetContract(contractHash string) (core.Contract, error)
	ContractExists(contractHash string) bool
}

var ErrFilterDoesNotExist = func(name string) error {
	return fmt.Errorf("filter %s does not exist", name)
}

type FilterRepository interface {
	CreateFilter(filter filters.LogFilter) error
	GetFilter(name string) (filters.LogFilter, error)
}

type HeaderRepository interface {
	CreateOrUpdateHeader(header core.Header) (int64, error)
	GetHeader(blockNumber int64) (core.Header, error)
	MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64
}

type LogRepository interface {
	CreateLogs(logs []core.Log, receiptId int64) error
	GetLogs(address string, blockNumber int64) []core.Log
}

var ErrReceiptDoesNotExist = func(txHash string) error {
	return fmt.Errorf("Receipt for tx: %v does not exist", txHash)
}

type ReceiptRepository interface {
	CreateReceiptsAndLogs(blockId int64, receipts []core.Receipt) error
	CreateReceipt(blockId int64, receipt core.Receipt) (int64, error)
	GetReceipt(txHash string) (core.Receipt, error)
}

type WatchedEventRepository interface {
	GetWatchedEvents(name string) ([]*core.WatchedEvent, error)
}
