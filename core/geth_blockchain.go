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
func (blockchain GethBlockchain) notifyObservers(header *types.Header) {
	block := Block{Number: header.Number}
	for _, observer := range blockchain.observers {
		observer.NotifyBlockAdded(block)
	}
}

func (blockchain *GethBlockchain) RegisterObserver(observer BlockchainObserver) {
	fmt.Printf("Registering observer: %v\n", reflect.TypeOf(observer))
	blockchain.observers = append(blockchain.observers, observer)
}

func (blockchain *GethBlockchain) SubscribeToEvents() {
	blocks := make(chan *types.Header, 10)
	myContext := context.Background()
	blockchain.client.SubscribeNewHead(myContext, blocks)
	for block := range blocks {
		blockchain.notifyObservers(block)
	}
}
