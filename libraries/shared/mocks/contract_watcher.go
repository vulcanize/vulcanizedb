package mocks

import (
	"fmt"
	"sync"
)

type MockContractWatcher struct {
	AddTransformersWasCalled  bool
	ExecuteWasCalled          bool
	WatchEthContractWasCalled bool
}

func NewMockContractWatcher() MockContractWatcher {
	return MockContractWatcher{
		AddTransformersWasCalled:  false,
		ExecuteWasCalled:          false,
		WatchEthContractWasCalled: false,
	}
}

func (cw *MockContractWatcher) AddTransformers(inits interface{}) error {
	cw.AddTransformersWasCalled = true
	return nil
}

func (cw *MockContractWatcher) Execute() error {
	cw.ExecuteWasCalled = true
	return nil
}

func (cw *MockContractWatcher) WatchEthContract(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("fake watch eth contract")
	cw.WatchEthContractWasCalled = true
}
