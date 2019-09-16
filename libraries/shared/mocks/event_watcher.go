package mocks

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"sync"
)

type MockEventWatcher struct {
	AddTransformersWasCalled bool
	ExecuteWasCalled         bool
	WatchEthEventsWasCalled  bool
	Transformers             []transformer.EventTransformer
}

func NewMockEventWatcher() MockEventWatcher {
	return MockEventWatcher{
		AddTransformersWasCalled: false,
		ExecuteWasCalled:         false,
		WatchEthEventsWasCalled:  false,
	}
}

func (ew *MockEventWatcher) Execute() error {
	ew.ExecuteWasCalled = true
	return nil
}

func (ew *MockEventWatcher) AddTransformers(initializers []transformer.EventTransformerInitializer) {
	fakeTransformer := &MockTransformer{}
	fakeTransformer.SetTransformerConfig(FakeTransformerConfig)
	ew.Transformers = []transformer.EventTransformer{fakeTransformer}
	ew.AddTransformersWasCalled = true
}

func (ew *MockEventWatcher) WatchEthEvents(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("fake watch eth events")
	ew.WatchEthEventsWasCalled = true
}
