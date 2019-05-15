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

package core

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockChain interface {
	ContractDataFetcher
	AccountDataFetcher
	GetBlockByNumber(blockNumber int64) (Block, error)
	GetEthLogsWithCustomQuery(query ethereum.FilterQuery) ([]types.Log, error)
	GetHeaderByNumber(blockNumber int64) (Header, error)
	GetHeadersByNumbers(blockNumbers []int64) ([]Header, error)
	GetLogs(contract Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]Log, error)
	GetTransactions(transactionHashes []common.Hash) ([]TransactionModel, error)
	LastBlock() (*big.Int, error)
	Node() Node
}

type ContractDataFetcher interface {
	FetchContractData(abiJSON string, address string, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error
}

type AccountDataFetcher interface {
	GetAccountBalance(address common.Address, blockNumber *big.Int) (*big.Int, error)
}
