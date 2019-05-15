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
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockChain struct {
	fetchContractDataErr               error
	fetchContractDataPassedAbi         string
	fetchContractDataPassedAddress     string
	fetchContractDataPassedMethod      string
	fetchContractDataPassedMethodArgs  []interface{}
	fetchContractDataPassedResult      interface{}
	fetchContractDataPassedBlockNumber int64
	getBlockByNumberErr                error
	GetTransactionsCalled              bool
	GetTransactionsError               error
	GetTransactionsPassedHashes        []common.Hash
	logQuery                           ethereum.FilterQuery
	logQueryErr                        error
	logQueryReturnLogs                 []types.Log
	lastBlock                          *big.Int
	node                               core.Node
	Transactions                       []core.TransactionModel
	accountBalanceReturnValue          *big.Int
	getAccountBalanceErr               error
}

func NewMockBlockChain() *MockBlockChain {
	return &MockBlockChain{
		node: core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "Geth"},
	}
}

func (chain *MockBlockChain) SetFetchContractDataErr(err error) {
	chain.fetchContractDataErr = err
}

func (chain *MockBlockChain) SetLastBlock(blockNumber *big.Int) {
	chain.lastBlock = blockNumber
}

func (chain *MockBlockChain) SetGetBlockByNumberErr(err error) {
	chain.getBlockByNumberErr = err
}

func (chain *MockBlockChain) SetGetEthLogsWithCustomQueryErr(err error) {
	chain.logQueryErr = err
}

func (chain *MockBlockChain) SetGetEthLogsWithCustomQueryReturnLogs(logs []types.Log) {
	chain.logQueryReturnLogs = logs
}

func (blockChain *MockBlockChain) FetchContractData(abiJSON string, address string, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	blockChain.fetchContractDataPassedAbi = abiJSON
	blockChain.fetchContractDataPassedAddress = address
	blockChain.fetchContractDataPassedMethod = method
	blockChain.fetchContractDataPassedMethodArgs = methodArgs
	blockChain.fetchContractDataPassedResult = result
	blockChain.fetchContractDataPassedBlockNumber = blockNumber
	return blockChain.fetchContractDataErr
}

func (chain *MockBlockChain) GetBlockByNumber(blockNumber int64) (core.Block, error) {
	return core.Block{Number: blockNumber}, chain.getBlockByNumberErr
}

func (blockChain *MockBlockChain) GetEthLogsWithCustomQuery(query ethereum.FilterQuery) ([]types.Log, error) {
	blockChain.logQuery = query
	return blockChain.logQueryReturnLogs, blockChain.logQueryErr
}

func (chain *MockBlockChain) GetHeaderByNumber(blockNumber int64) (core.Header, error) {
	return core.Header{BlockNumber: blockNumber}, nil
}

func (chain *MockBlockChain) GetHeadersByNumbers(blockNumbers []int64) ([]core.Header, error) {
	var headers []core.Header
	for _, blockNumber := range blockNumbers {
		var header = core.Header{BlockNumber: int64(blockNumber)}
		headers = append(headers, header)
	}
	return headers, nil
}

func (chain *MockBlockChain) GetLogs(contract core.Contract, startingBlockNumber, endingBlockNumber *big.Int) ([]core.Log, error) {
	return []core.Log{}, nil
}

func (chain *MockBlockChain) GetTransactions(transactionHashes []common.Hash) ([]core.TransactionModel, error) {
	chain.GetTransactionsCalled = true
	chain.GetTransactionsPassedHashes = transactionHashes
	return chain.Transactions, chain.GetTransactionsError
}

func (chain *MockBlockChain) CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (chain *MockBlockChain) LastBlock() (*big.Int, error) {
	return chain.lastBlock, nil
}

func (chain *MockBlockChain) Node() core.Node {
	return chain.node
}

func (chain *MockBlockChain) AssertFetchContractDataCalledWith(abiJSON string, address string, method string, methodArgs []interface{}, result interface{}, blockNumber int64) {
	Expect(chain.fetchContractDataPassedAbi).To(Equal(abiJSON))
	Expect(chain.fetchContractDataPassedAddress).To(Equal(address))
	Expect(chain.fetchContractDataPassedMethod).To(Equal(method))
	if methodArgs != nil {
		Expect(chain.fetchContractDataPassedMethodArgs).To(Equal(methodArgs))
	}
	Expect(chain.fetchContractDataPassedResult).To(BeAssignableToTypeOf(result))
	Expect(chain.fetchContractDataPassedBlockNumber).To(Equal(blockNumber))
}

func (blockChain *MockBlockChain) AssertGetEthLogsWithCustomQueryCalledWith(query ethereum.FilterQuery) {
	Expect(blockChain.logQuery).To(Equal(query))
}

func (blockChain *MockBlockChain) SetGetAccountBalanceErr(err error) {
	blockChain.getAccountBalanceErr = err
}

func (blockChain *MockBlockChain) SetGetAccountBalance(balance *big.Int) {
	blockChain.accountBalanceReturnValue = balance
}

func (blockChain *MockBlockChain) GetAccountBalance(address common.Address, blockNumber *big.Int) (*big.Int, error) {
	return blockChain.accountBalanceReturnValue, blockChain.getAccountBalanceErr
}
