package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"
	"reflect"
)

type GethBlockchain struct {
	client    *ethclient.Client
	observers []BlockchainObserver
}

func NewGethBlockchain(ipcPath string) *GethBlockchain {
	fmt.Printf("Creating Geth Blockchain to: %s\n", ipcPath)
	blockchain := GethBlockchain{}
	client, _ := ethclient.Dial(ipcPath)
	// TODO: handle error gracefully
	blockchain.client = client
	return &blockchain
}
func (blockchain GethBlockchain) notifyObservers(getBlock *types.Block) {
	block := convertBlock(getBlock)
	for _, observer := range blockchain.observers {
		observer.NotifyBlockAdded(block)
	}
}
func convertBlock(gethBlock *types.Block) Block {
	return Block{
		Number:               gethBlock.Number(),
		NumberOfTransactions: len(gethBlock.Transactions()),
	}
}

func (blockchain *GethBlockchain) RegisterObserver(observer BlockchainObserver) {
	fmt.Printf("Registering observer: %v\n", reflect.TypeOf(observer))
	blockchain.observers = append(blockchain.observers, observer)
}

func (blockchain *GethBlockchain) SubscribeToEvents() {
	headers := make(chan *types.Header, 10)
	myContext := context.Background()
	blockchain.client.SubscribeNewHead(myContext, headers)
	for header := range headers {
		gethBlock, _ := blockchain.client.BlockByNumber(myContext, header.Number)
		blockchain.notifyObservers(gethBlock)
	}
}
