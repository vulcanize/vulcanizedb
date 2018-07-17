package geth

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
	vulcRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

type BlockChain struct {
	client          *ethclient.Client
	blockConverter  vulcCommon.BlockConverter
	headerConverter vulcCommon.HeaderConverter
	node            core.Node
}

func NewBlockChain(ipcPath string) *BlockChain {
	rpcClient, err := rpc.Dial(ipcPath)
	if err != nil {
		log.Fatal(err)
	}
	client := ethclient.NewClient(rpcClient)
	clientWrapper := node.ClientWrapper{ContextCaller: rpcClient, IPCPath: ipcPath}
	transactionConverter := vulcRpc.NewRpcTransactionConverter(client)
	return &BlockChain{
		client:          client,
		blockConverter:  vulcCommon.NewBlockConverter(transactionConverter),
		headerConverter: vulcCommon.HeaderConverter{},
		node:            node.MakeNode(clientWrapper),
	}
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
	}
	gethLogs, err := blockChain.client.FilterLogs(context.Background(), fc)
	if err != nil {
		return []core.Log{}, err
	}
	logs := vulcCommon.ToCoreLogs(gethLogs)
	return logs, nil
}

func (blockChain *BlockChain) Node() core.Node {
	return blockChain.node
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

func (blockChain *BlockChain) LastBlock() *big.Int {
	block, _ := blockChain.client.HeaderByNumber(context.Background(), nil)
	return block.Number
}
