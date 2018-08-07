package fakes

import (
	"math/big"

	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
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
	lastBlock                          *big.Int
	node                               core.Node
}

func NewMockBlockChain() *MockBlockChain {
	return &MockBlockChain{
		node: core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "Geth"},
	}
}

func (blockChain *MockBlockChain) SetFetchContractDataErr(err error) {
	blockChain.fetchContractDataErr = err
}

func (blockChain *MockBlockChain) SetLastBlock(blockNumber *big.Int) {
	blockChain.lastBlock = blockNumber
}

func (blockChain *MockBlockChain) SetGetBlockByNumberErr(err error) {
	blockChain.getBlockByNumberErr = err
}

func (blockChain *MockBlockChain) SetGetLogsErr(err error) {
	blockChain.logQueryErr = err
}

func (blockChain *MockBlockChain) GetEthLogsWithCustomQuery(query ethereum.FilterQuery) ([]types.Log, error) {
	blockChain.logQuery = query
	return []types.Log{}, blockChain.logQueryErr
}

func (blockChain *MockBlockChain) GetHeaderByNumber(blockNumber int64) (core.Header, error) {
	return core.Header{BlockNumber: blockNumber}, nil
}

func (blockChain *MockBlockChain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	blockChain.fetchContractDataPassedAbi = abiJSON
	blockChain.fetchContractDataPassedAddress = address
	blockChain.fetchContractDataPassedMethod = method
	blockChain.fetchContractDataPassedMethodArg = methodArg
	blockChain.fetchContractDataPassedResult = result
	blockChain.fetchContractDataPassedBlockNumber = blockNumber
	return blockChain.fetchContractDataErr
}

func (blockChain *MockBlockChain) CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (blockChain *MockBlockChain) LastBlock() *big.Int {
	return blockChain.lastBlock
}

func (blockChain *MockBlockChain) GetLogs(contract core.Contract, startingBlock *big.Int, endingBlock *big.Int) ([]core.Log, error) {
	return []core.Log{}, nil
}

func (blockChain *MockBlockChain) Node() core.Node {
	return blockChain.node
}

func (blockChain *MockBlockChain) GetBlockByNumber(blockNumber int64) (core.Block, error) {
	return core.Block{Number: blockNumber}, blockChain.getBlockByNumberErr
}

// TODO: handle methodArg being nil (can't match nil to nil in Gomega)
func (blockChain *MockBlockChain) AssertFetchContractDataCalledWith(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) {
	Expect(blockChain.fetchContractDataPassedAbi).To(Equal(abiJSON))
	Expect(blockChain.fetchContractDataPassedAddress).To(Equal(address))
	Expect(blockChain.fetchContractDataPassedMethod).To(Equal(method))
	Expect(blockChain.fetchContractDataPassedResult).To(Equal(result))
	Expect(blockChain.fetchContractDataPassedBlockNumber).To(Equal(blockNumber))
}

func (blockChain *MockBlockChain) AssertGetEthLogsWithCustomQueryCalledWith(query ethereum.FilterQuery) {
	Expect(blockChain.logQuery).To(Equal(query))
}
