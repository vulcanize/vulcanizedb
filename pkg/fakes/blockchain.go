package fakes

import (
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Blockchain struct {
	logs                map[string][]core.Log
	blocks              map[int64]core.Block
	contractAttributes  map[string]map[string]string
	blocksChannel       chan core.Block
	WasToldToStop       bool
	node                core.Node
	ContractReturnValue []byte
}

func (blockchain *Blockchain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	panic("implement me")
}

func (blockchain *Blockchain) CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error) {
	return blockchain.ContractReturnValue, nil
}

func (blockchain *Blockchain) LastBlock() *big.Int {
	var max int64
	for blockNumber := range blockchain.blocks {
		if blockNumber > max {
			max = blockNumber
		}
	}
	return big.NewInt(max)
}

func (blockchain *Blockchain) GetLogs(contract core.Contract, startingBlock *big.Int, endingBlock *big.Int) ([]core.Log, error) {
	return blockchain.logs[contract.Hash], nil
}

func (blockchain *Blockchain) Node() core.Node {
	return blockchain.node
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks:             make(map[int64]core.Block),
		logs:               make(map[string][]core.Log),
		contractAttributes: make(map[string]map[string]string),
		node:               core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "Geth"},
	}
}

func NewBlockchainWithBlocks(blocks []core.Block) *Blockchain {
	blockNumberToBlocks := make(map[int64]core.Block)
	for _, block := range blocks {
		blockNumberToBlocks[block.Number] = block
	}
	return &Blockchain{
		blocks: blockNumberToBlocks,
	}
}

func (blockchain *Blockchain) GetBlockByNumber(blockNumber int64) core.Block {
	return blockchain.blocks[blockNumber]
}

func (blockchain *Blockchain) AddBlock(block core.Block) {
	blockchain.blocks[block.Number] = block
	blockchain.blocksChannel <- block
}
