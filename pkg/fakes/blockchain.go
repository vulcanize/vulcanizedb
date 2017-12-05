package fakes

import (
	"sort"

	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/core"
)

type Blockchain struct {
	blocks             map[int64]core.Block
	contractAttributes map[string]map[string]string
	blocksChannel      chan core.Block
	WasToldToStop      bool
}

func (blockchain *Blockchain) GetAttribute(contract core.Contract, attributeName string, blockNumber *big.Int) (interface{}, error) {
	var result interface{}
	if blockNumber == nil {
		result = blockchain.contractAttributes[contract.Hash+"-1"][attributeName]
	} else {
		result = blockchain.contractAttributes[contract.Hash+blockNumber.String()][attributeName]
	}
	return result, nil
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks:             make(map[int64]core.Block),
		contractAttributes: make(map[string]map[string]string),
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

func (blockchain *Blockchain) SubscribeToBlocks(outputBlocks chan core.Block) {
	blockchain.blocksChannel = outputBlocks
}

func (blockchain *Blockchain) AddBlock(block core.Block) {
	blockchain.blocks[block.Number] = block
	blockchain.blocksChannel <- block
}

func (*Blockchain) StartListening() {}

func (blockchain *Blockchain) StopListening() {
	blockchain.WasToldToStop = true
}

func (blockchain *Blockchain) SetContractStateAttribute(contractHash string, blockNumber *big.Int, attributeName string, attributeValue string) {
	var key string
	if blockNumber == nil {
		key = contractHash + "-1"
	} else {
		key = contractHash + blockNumber.String()
	}
	contractStateAttributes := blockchain.contractAttributes[key]
	if contractStateAttributes == nil {
		blockchain.contractAttributes[key] = make(map[string]string)
	}
	blockchain.contractAttributes[key][attributeName] = attributeValue
}

func (blockchain *Blockchain) GetAttributes(contract core.Contract) (core.ContractAttributes, error) {
	var contractAttributes core.ContractAttributes
	attributes, ok := blockchain.contractAttributes[contract.Hash+"-1"]
	if ok {
		for key, _ := range attributes {
			contractAttributes = append(contractAttributes, core.ContractAttribute{Name: key, Type: "string"})
		}
	}
	sort.Sort(contractAttributes)
	return contractAttributes, nil
}
