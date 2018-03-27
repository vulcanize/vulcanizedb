package geth

import (
	"math/big"

	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"golang.org/x/net/context"
)

type Blockchain struct {
	client              *ethclient.Client
	readGethHeaders     chan *types.Header
	outputBlocks        chan core.Block
	newHeadSubscription ethereum.Subscription
	node                core.Node
}

func NewBlockchain(ipcPath string) *Blockchain {
	blockchain := Blockchain{}
	rpcClient, err := rpc.Dial(ipcPath)
	if err != nil {
		log.Fatal(err)
	}
	client := ethclient.NewClient(rpcClient)
	clientWrapper := node.ClientWrapper{ContextCaller: rpcClient, IPCPath: ipcPath}
	blockchain.node = node.MakeNode(clientWrapper)
	blockchain.client = client
	return &blockchain
}

func (blockchain *Blockchain) GetLogs(contract core.Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]core.Log, error) {
	if endingBlockNumber == nil {
		endingBlockNumber = startingBlockNumber
	}
	contractAddress := common.HexToAddress(contract.Hash)
	fc := ethereum.FilterQuery{
		FromBlock: startingBlockNumber,
		ToBlock:   endingBlockNumber,
		Addresses: []common.Address{contractAddress},
	}
	gethLogs, err := blockchain.client.FilterLogs(context.Background(), fc)
	if err != nil {
		return []core.Log{}, err
	}
	logs := ToCoreLogs(gethLogs)
	return logs, nil
}

func (blockchain *Blockchain) Node() core.Node {
	return blockchain.node
}

func (blockchain *Blockchain) GetBlockByNumber(blockNumber int64) (core.Block, error) {
	gethBlock, err := blockchain.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		return core.Block{}, err
	}
	block, err := ToCoreBlock(gethBlock, blockchain.client)
	if err != nil {
		return core.Block{}, err
	}
	return block, nil
}

func (blockchain *Blockchain) LastBlock() *big.Int {
	block, _ := blockchain.client.HeaderByNumber(context.Background(), nil)
	return block.Number
}
