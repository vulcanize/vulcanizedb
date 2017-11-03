package geth

import (
	"fmt"

	"github.com/8thlight/vulcanizedb/core"
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
	myContext := context.Background()
	for header := range blockchain.readGethHeaders {
		gethBlock, _ := blockchain.client.BlockByNumber(myContext, header.Number)
		block := GethBlockToCoreBlock(gethBlock)
		blockchain.outputBlocks <- block
	}
}

func (blockchain *GethBlockchain) StopListening() {
	blockchain.newHeadSubscription.Unsubscribe()
}
