package core

import (
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	. "github.com/onsi/gomega"
)

type MockBlockChain struct {
	configCalled bool
}

func NewMockBlockChain() *MockBlockChain {
	return &MockBlockChain{
		configCalled: false,
	}
}

func (*MockBlockChain) BlockChain() *core.BlockChain {
	panic("implement me")
}

func (mbc *MockBlockChain) Config() *params.ChainConfig {
	mbc.configCalled = true
	return params.TestChainConfig
}

func (*MockBlockChain) Engine() consensus.Engine {
	panic("implement me")
}

func (mbc *MockBlockChain) AssertConfigCalled() {
	Expect(mbc.configCalled).To(BeTrue())
}
