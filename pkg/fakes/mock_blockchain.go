package fakes

import (
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
	fetchContractDataPassedMethodArg   interface{}
	fetchContractDataPassedResult      interface{}
	fetchContractDataPassedBlockNumber int64
	getBlockByNumberErr                error
	logQuery                           ethereum.FilterQuery
	logQueryErr                        error
	logQueryReturnLogs                 []types.Log
	lastBlock                          *big.Int
	node                               core.Node
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

func (chain *MockBlockChain) FetchContractData(abiJSON, address, method string, methodArg, result interface{}, blockNumber int64) error {
	chain.fetchContractDataPassedAbi = abiJSON
	chain.fetchContractDataPassedAddress = address
	chain.fetchContractDataPassedMethod = method
	chain.fetchContractDataPassedMethodArg = methodArg
	chain.fetchContractDataPassedResult = result
	chain.fetchContractDataPassedBlockNumber = blockNumber
	return chain.fetchContractDataErr
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

func (chain *MockBlockChain) GetHeaderByNumbers(blockNumbers []int64) ([]core.Header, error) {
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

func (chain *MockBlockChain) CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (chain *MockBlockChain) LastBlock() *big.Int {
	return chain.lastBlock
}

func (chain *MockBlockChain) Node() core.Node {
	return chain.node
}

func (chain *MockBlockChain) AssertFetchContractDataCalledWith(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) {
	Expect(chain.fetchContractDataPassedAbi).To(Equal(abiJSON))
	Expect(chain.fetchContractDataPassedAddress).To(Equal(address))
	Expect(chain.fetchContractDataPassedMethod).To(Equal(method))
	if methodArg != nil {
		Expect(chain.fetchContractDataPassedMethodArg).To(Equal(methodArg))
	}
	Expect(chain.fetchContractDataPassedResult).To(BeAssignableToTypeOf(result))
	Expect(chain.fetchContractDataPassedBlockNumber).To(Equal(blockNumber))
}

func (blockChain *MockBlockChain) AssertGetEthLogsWithCustomQueryCalledWith(query ethereum.FilterQuery) {
	Expect(blockChain.logQuery).To(Equal(query))
}
