package core

import "fmt"

type GethBlockchain struct{}

func NewGethBlockchain() *GethBlockchain {
	fmt.Println("Creating Gethblockchain")
	return &GethBlockchain{}
}

func (blockchain *GethBlockchain) RegisterObserver(_ BlockchainObserver) {

}
