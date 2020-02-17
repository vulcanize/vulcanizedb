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
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"

	. "github.com/onsi/gomega"
)

type MockEthClient struct {
	callContractErr             error
	callContractPassedContext   context.Context
	callContractPassedMsg       ethereum.CallMsg
	callContractPassedNumber    *big.Int
	callContractReturnBytes     []byte
	blockByNumberErr            error
	blockByNumberPassedContext  context.Context
	blockByNumberPassedNumber   *big.Int
	blockByNumberReturnBlock    *types.Block
	headerByNumberErr           error
	headerByNumberPassedContext context.Context
	headerByNumberPassedNumber  *big.Int
	headerByNumberReturnHeader  *types.Header
	headerByNumbersReturnHeader []*types.Header
	headerByNumbersPassedNumber []*big.Int
	filterLogsErr               error
	filterLogsPassedContext     context.Context
	filterLogsPassedQuery       ethereum.FilterQuery
	filterLogsReturnLogs        []types.Log
	transactionReceipts         map[string]*types.Receipt
	storageAtCtx                context.Context
	storageAtAccount            common.Address
	storageAtKey                common.Hash
	storageAtBlock              *big.Int
	storageAtError              error
	err                         error
	passedBatch                 []client.BatchElem
	passedMethod                string
	transactionSenderErr        error
	transactionReceiptErr       error
	passedAddress               common.Address
	passedBlockNumber           *big.Int
	balanceToReturn             *big.Int
	balanceAtErr                error
	passedbalanceAtContext      context.Context
}

func NewMockEthClient() *MockEthClient {
	return &MockEthClient{}
}

func (client *MockEthClient) SetCallContractErr(err error) {
	client.callContractErr = err
}

func (client *MockEthClient) SetCallContractReturnBytes(returnBytes []byte) {
	client.callContractReturnBytes = returnBytes
}

func (client *MockEthClient) SetBlockByNumberErr(err error) {
	client.blockByNumberErr = err
}

func (client *MockEthClient) SetBlockByNumberReturnBlock(block *types.Block) {
	client.blockByNumberReturnBlock = block
}

func (client *MockEthClient) SetHeaderByNumberErr(err error) {
	client.headerByNumberErr = err
}

func (client *MockEthClient) SetHeaderByNumberReturnHeader(header *types.Header) {
	client.headerByNumberReturnHeader = header
}

func (client *MockEthClient) SetHeaderByNumbersReturnHeader(headers []*types.Header) {
	client.headerByNumbersReturnHeader = headers
}

func (client *MockEthClient) SetFilterLogsErr(err error) {
	client.filterLogsErr = err
}

func (client *MockEthClient) SetFilterLogsReturnLogs(logs []types.Log) {
	client.filterLogsReturnLogs = logs
}

func (client *MockEthClient) SetTransactionReceiptErr(err error) {
	client.transactionReceiptErr = err
}

func (client *MockEthClient) SetTransactionReceipts(receipts []*types.Receipt) {
	for _, receipt := range receipts {
		client.transactionReceipts[receipt.TxHash.Hex()] = receipt
	}
}

func (client *MockEthClient) SetTransactionSenderErr(err error) {
	client.transactionSenderErr = err
}

func (client *MockEthClient) SetBalanceAtErr(err error) {
	client.balanceAtErr = err
}

func (client *MockEthClient) SetBalanceAt(balance *big.Int) {
	client.balanceToReturn = balance
}

func (client *MockEthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	client.callContractPassedContext = ctx
	client.callContractPassedMsg = msg
	client.callContractPassedNumber = blockNumber
	return client.callContractReturnBytes, client.callContractErr
}

func (client *MockEthClient) BatchCall(batch []client.BatchElem) error {
	client.passedBatch = batch
	client.passedMethod = batch[0].Method

	return nil
}

func (client *MockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	client.blockByNumberPassedContext = ctx
	client.blockByNumberPassedNumber = number
	return client.blockByNumberReturnBlock, client.blockByNumberErr
}

func (client *MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	client.headerByNumberPassedContext = ctx
	client.headerByNumberPassedNumber = number
	return client.headerByNumberReturnHeader, client.headerByNumberErr
}

func (client *MockEthClient) HeaderByNumbers(numbers []*big.Int) ([]*types.Header, error) {
	client.headerByNumbersPassedNumber = numbers
	return client.headerByNumbersReturnHeader, client.headerByNumberErr
}

func (client *MockEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	client.filterLogsPassedContext = ctx
	client.filterLogsPassedQuery = q
	return client.filterLogsReturnLogs, client.filterLogsErr
}

func (client *MockEthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	return common.HexToAddress("0x123"), client.transactionSenderErr
}

func (client *MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if gasUsed, ok := client.transactionReceipts[txHash.Hex()]; ok {
		return gasUsed, client.transactionReceiptErr
	}
	return &types.Receipt{GasUsed: uint64(0)}, client.transactionReceiptErr
}

func (client *MockEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	client.passedbalanceAtContext = ctx
	client.passedAddress = account
	client.passedBlockNumber = blockNumber
	return client.balanceToReturn, client.balanceAtErr
}

func (client *MockEthClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	client.storageAtCtx = ctx
	client.storageAtAccount = account
	client.storageAtKey = key
	client.storageAtBlock = blockNumber
	return []byte{}, client.storageAtError
}

func (client *MockEthClient) SetStorageAtError(err error) {
	client.storageAtError = err
}

func (client *MockEthClient) AssertStorageAtCalledWith(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) {
	Expect(client.storageAtCtx).To(Equal(ctx))
	Expect(client.storageAtAccount).To(Equal(account))
	Expect(client.storageAtKey).To(Equal(key))
	Expect(client.storageAtBlock).To(Equal(blockNumber))
}

func (client *MockEthClient) AssertCallContractCalledWith(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) {
	Expect(client.callContractPassedContext).To(Equal(ctx))
	Expect(client.callContractPassedMsg).To(Equal(msg))
	Expect(client.callContractPassedNumber).To(Equal(blockNumber))
}

func (client *MockEthClient) AssertBlockByNumberCalledWith(ctx context.Context, number *big.Int) {
	Expect(client.blockByNumberPassedContext).To(Equal(ctx))
	Expect(client.blockByNumberPassedNumber).To(Equal(number))
}

func (client *MockEthClient) AssertHeaderByNumberCalledWith(ctx context.Context, number *big.Int) {
	Expect(client.headerByNumberPassedContext).To(Equal(ctx))
	Expect(client.headerByNumberPassedNumber).To(Equal(number))
}

func (client *MockEthClient) AssertHeaderByNumbersCalledWith(number []*big.Int) {
	Expect(client.headerByNumbersPassedNumber).To(Equal(number))
}

func (client *MockEthClient) AssertFilterLogsCalledWith(ctx context.Context, q ethereum.FilterQuery) {
	Expect(client.filterLogsPassedContext).To(Equal(ctx))
	Expect(client.filterLogsPassedQuery).To(Equal(q))
}

func (client *MockEthClient) AssertBatchCalledWith(method string) {
	Expect(client.passedMethod).To(Equal(method))
}

func (client *MockEthClient) AssertBalanceAtCalled(ctx context.Context, account common.Address, blockNumber *big.Int) {
	Expect(client.passedbalanceAtContext).To(Equal(ctx))
	Expect(client.passedAddress).To(Equal(account))
	Expect(client.passedBlockNumber).To(Equal(blockNumber))
}
