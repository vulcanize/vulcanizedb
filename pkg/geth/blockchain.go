package geth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/net/context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

type BlockChain struct {
	client          core.EthClient
	blockConverter  vulcCommon.BlockConverter
	headerConverter vulcCommon.HeaderConverter
	node            core.Node
}

func NewBlockChain(client core.EthClient, node core.Node, converter vulcCommon.TransactionConverter) *BlockChain {
	return &BlockChain{
		client:          client,
		blockConverter:  vulcCommon.NewBlockConverter(converter),
		headerConverter: vulcCommon.HeaderConverter{},
		node:            node,
	}
}

func (blockChain *BlockChain) GetBlockByNumber(blockNumber int64) (block core.Block, err error) {
	gethBlock, err := blockChain.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return block, err
	}
	return blockChain.blockConverter.ToCoreBlock(gethBlock)
}

func (blockChain *BlockChain) GetHeaderByNumber(blockNumber int64) (header core.Header, err error) {
	gethHeader, err := blockChain.client.HeaderByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return header, err
	}
	return blockChain.headerConverter.Convert(gethHeader)
}

func (blockChain *BlockChain) GetLogs(contract core.Contract, startingBlockNumber, endingBlockNumber *big.Int) ([]core.Log, error) {
	if endingBlockNumber == nil {
		endingBlockNumber = startingBlockNumber
	}
	contractAddress := common.HexToAddress(contract.Hash)
	fc := ethereum.FilterQuery{
		FromBlock: startingBlockNumber,
		ToBlock:   endingBlockNumber,
		Addresses: []common.Address{contractAddress},
		Topics:    nil,
	}
	gethLogs, err := blockChain.GetEthLogsWithCustomQuery(fc)
	if err != nil {
		return []core.Log{}, err
	}
	logs := vulcCommon.ToCoreLogs(gethLogs)
	return logs, nil
}

func (blockChain *BlockChain) GetEthLogsWithCustomQuery(query ethereum.FilterQuery) ([]types.Log, error) {
	gethLogs, err := blockChain.client.FilterLogs(context.Background(), query)
	if err != nil {
		return []types.Log{}, err
	}
	return gethLogs, nil
}

func (blockChain *BlockChain) LastBlock() *big.Int {
	block, _ := blockChain.client.HeaderByNumber(context.Background(), nil)
	return block.Number
}

func (blockChain *BlockChain) Node() core.Node {
	return blockChain.node
}
