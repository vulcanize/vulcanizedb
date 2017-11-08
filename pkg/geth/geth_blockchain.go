package geth

import (
	"fmt"

	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"
)

type GethBlockchain struct {
	client              *ethclient.Client
	readGethHeaders     chan *types.Header
	outputBlocks        chan core.Block
	newHeadSubscription ethereum.Subscription
}

func (blockchain *GethBlockchain) GetBlockByNumber(blockNumber int64) core.Block {
	gethBlock, _ := blockchain.client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	return GethBlockToCoreBlock(gethBlock)
}

func NewGethBlockchain(ipcPath string) *GethBlockchain {
	fmt.Printf("Creating Geth Blockchain to: %s\n", ipcPath)
	blockchain := GethBlockchain{}
	client, _ := ethclient.Dial(ipcPath)
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
