package geth

import (
	"fmt"

	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/geth/node"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
)

type GethBlockchain struct {
	client              *ethclient.Client
	readGethHeaders     chan *types.Header
	outputBlocks        chan core.Block
	newHeadSubscription ethereum.Subscription
	node                core.Node
}

func (blockchain *GethBlockchain) GetLogs(contract core.Contract, blockNumber *big.Int) ([]core.Log, error) {
	if blockNumber == nil {
		blockNumber = blockchain.latestBlock()
	}
	contractAddress := common.HexToAddress(contract.Hash)
	fc := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Addresses: []common.Address{contractAddress},
	}
	gethLogs, err := blockchain.client.FilterLogs(context.Background(), fc)
	if err != nil {
		return []core.Log{}, err
	}
	logs := GethLogsToCoreLogs(gethLogs)
	return logs, nil
}

func (blockchain *GethBlockchain) Node() core.Node {
	return blockchain.node
}

func (blockchain *GethBlockchain) GetBlockByNumber(blockNumber int64) core.Block {
	gethBlock, _ := blockchain.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	return GethBlockToCoreBlock(gethBlock, blockchain.client)
}

func NewGethBlockchain(ipcPath string) *GethBlockchain {
	blockchain := GethBlockchain{}
	rpcClient, _ := rpc.Dial(ipcPath)
	client := ethclient.NewClient(rpcClient)
	blockchain.node = node.Retrieve(rpcClient)
	blockchain.client = client
	return &blockchain
}

func (blockchain *GethBlockchain) SubscribeToBlocks(blocks chan core.Block) {
	blockchain.outputBlocks = blocks
	fmt.Println("SubscribeToBlocks")
	inputHeaders := make(chan *types.Header, 10)
	myContext := context.Background()
	blockchain.readGethHeaders = inputHeaders
	subscription, _ := blockchain.client.SubscribeNewHead(myContext, inputHeaders)
	blockchain.newHeadSubscription = subscription
}

func (blockchain *GethBlockchain) StartListening() {
	for header := range blockchain.readGethHeaders {
		block := blockchain.GetBlockByNumber(header.Number.Int64())
		blockchain.outputBlocks <- block
	}
}

func (blockchain *GethBlockchain) StopListening() {
	blockchain.newHeadSubscription.Unsubscribe()
}

func (blockchain *GethBlockchain) latestBlock() *big.Int {
	block, _ := blockchain.client.HeaderByNumber(context.Background(), nil)
	return block.Number
}
