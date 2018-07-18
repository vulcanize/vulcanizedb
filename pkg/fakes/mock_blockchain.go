package fakes

import (
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockBlockChain struct {
	ContractReturnValue []byte
	WasToldToStop       bool
	blocks              map[int64]core.Block
	blocksChannel       chan core.Block
	contractAttributes  map[string]map[string]string
	err                 error
	headers             map[int64]core.Header
	logs                map[string][]core.Log
	node                core.Node
}

func (blockChain *MockBlockChain) GetHeaderByNumber(blockNumber int64) (core.Header, error) {
	return blockChain.headers[blockNumber], nil
}

func (blockChain *MockBlockChain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	panic("implement me")
}

func (blockChain *MockBlockChain) CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	return blockChain.ContractReturnValue, nil
}

func (blockChain *MockBlockChain) LastBlock() *big.Int {
	var max int64
	for blockNumber := range blockChain.blocks {
		if blockNumber > max {
			max = blockNumber
		}
	}
	return big.NewInt(max)
}

func (blockChain *MockBlockChain) GetLogs(contract core.Contract, startingBlock *big.Int, endingBlock *big.Int) ([]core.Log, error) {
	return blockChain.logs[contract.Hash], nil
}

func (blockChain *MockBlockChain) Node() core.Node {
	return blockChain.node
}

func NewMockBlockChain(err error) *MockBlockChain {
	return &MockBlockChain{
		blocks:             make(map[int64]core.Block),
		logs:               make(map[string][]core.Log),
		contractAttributes: make(map[string]map[string]string),
		node:               core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "Geth"},
		err:                err,
	}
}

func NewMockBlockChainWithBlocks(blocks []core.Block) *MockBlockChain {
	blockNumberToBlocks := make(map[int64]core.Block)
	for _, block := range blocks {
		blockNumberToBlocks[block.Number] = block
	}
	return &MockBlockChain{
		blocks: blockNumberToBlocks,
	}
}

func NewMockBlockChainWithHeaders(headers []core.Header) *MockBlockChain {
	// need to create blocks and headers so that LastBlock() will work in the mock
	// no reason to implement LastBlock() separately for headers since it checks
	// the last header in the Node's DB already
	memoryBlocks := make(map[int64]core.Block)
	memoryHeaders := make(map[int64]core.Header)
	for _, header := range headers {
		memoryBlocks[header.BlockNumber] = core.Block{Number: header.BlockNumber}
		memoryHeaders[header.BlockNumber] = header
	}
	return &MockBlockChain{
		blocks:  memoryBlocks,
		headers: memoryHeaders,
	}
}

func (blockChain *MockBlockChain) GetBlockByNumber(blockNumber int64) (core.Block, error) {
	if blockChain.err != nil {
		return core.Block{}, blockChain.err
	}
	return blockChain.blocks[blockNumber], nil
}

func (blockChain *MockBlockChain) AddBlock(block core.Block) {
	blockChain.blocks[block.Number] = block
	blockChain.blocksChannel <- block
}
