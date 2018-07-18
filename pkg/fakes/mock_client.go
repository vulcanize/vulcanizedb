package fakes

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockClient struct {
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
	filterLogsErr               error
	filterLogsPassedContext     context.Context
	filterLogsPassedQuery       ethereum.FilterQuery
	filterLogsReturnLogs        []types.Log
}

func NewMockClient() *MockClient {
	return &MockClient{
		callContractErr:             nil,
		callContractPassedContext:   nil,
		callContractPassedMsg:       ethereum.CallMsg{},
		callContractPassedNumber:    nil,
		callContractReturnBytes:     nil,
		blockByNumberErr:            nil,
		blockByNumberPassedContext:  nil,
		blockByNumberPassedNumber:   nil,
		blockByNumberReturnBlock:    nil,
		headerByNumberErr:           nil,
		headerByNumberPassedContext: nil,
		headerByNumberPassedNumber:  nil,
		headerByNumberReturnHeader:  nil,
		filterLogsErr:               nil,
		filterLogsPassedContext:     nil,
		filterLogsPassedQuery:       ethereum.FilterQuery{},
		filterLogsReturnLogs:        nil,
	}
}

func (client *MockClient) SetCallContractErr(err error) {
	client.callContractErr = err
}

func (client *MockClient) SetCallContractReturnBytes(returnBytes []byte) {
	client.callContractReturnBytes = returnBytes
}

func (client *MockClient) SetBlockByNumberErr(err error) {
	client.blockByNumberErr = err
}

func (client *MockClient) SetBlockByNumberReturnBlock(block *types.Block) {
	client.blockByNumberReturnBlock = block
}

func (client *MockClient) SetHeaderByNumberErr(err error) {
	client.headerByNumberErr = err
}

func (client *MockClient) SetHeaderByNumberReturnHeader(header *types.Header) {
	client.headerByNumberReturnHeader = header
}

func (client *MockClient) SetFilterLogsErr(err error) {
	client.filterLogsErr = err
}

func (client *MockClient) SetFilterLogsReturnLogs(logs []types.Log) {
	client.filterLogsReturnLogs = logs
}

func (client *MockClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	client.callContractPassedContext = ctx
	client.callContractPassedMsg = msg
	client.callContractPassedNumber = blockNumber
	return client.callContractReturnBytes, client.callContractErr
}

func (client *MockClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	client.blockByNumberPassedContext = ctx
	client.blockByNumberPassedNumber = number
	return client.blockByNumberReturnBlock, client.blockByNumberErr
}

func (client *MockClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	client.headerByNumberPassedContext = ctx
	client.headerByNumberPassedNumber = number
	return client.headerByNumberReturnHeader, client.headerByNumberErr
}

func (client *MockClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	client.filterLogsPassedContext = ctx
	client.filterLogsPassedQuery = q
	return client.filterLogsReturnLogs, client.filterLogsErr
}

func (client *MockClient) AssertCallContractCalledWith(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) {
	Expect(client.callContractPassedContext).To(Equal(ctx))
	Expect(client.callContractPassedMsg).To(Equal(msg))
	Expect(client.callContractPassedNumber).To(Equal(blockNumber))
}

func (client *MockClient) AssertBlockByNumberCalledWith(ctx context.Context, number *big.Int) {
	Expect(client.blockByNumberPassedContext).To(Equal(ctx))
	Expect(client.blockByNumberPassedNumber).To(Equal(number))
}

func (client *MockClient) AssertHeaderByNumberCalledWith(ctx context.Context, number *big.Int) {
	Expect(client.headerByNumberPassedContext).To(Equal(ctx))
	Expect(client.headerByNumberPassedNumber).To(Equal(number))
}

func (client *MockClient) AssertFilterLogsCalledWith(ctx context.Context, q ethereum.FilterQuery) {
	Expect(client.filterLogsPassedContext).To(Equal(ctx))
	Expect(client.filterLogsPassedQuery).To(Equal(q))
}
