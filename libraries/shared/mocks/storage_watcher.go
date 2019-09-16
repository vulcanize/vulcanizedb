package mocks

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"sync"
)

type MockStorageWatcher struct {
	AddTransformersWasCalled bool
	ExecuteWasCalled         bool
	WatchEthStorageWasCalled bool
	Transformers             []transformer.StorageTransformer
}

func NewMockStorageWatcher() MockStorageWatcher {
	return MockStorageWatcher{
		AddTransformersWasCalled: false,
		ExecuteWasCalled:         false,
		WatchEthStorageWasCalled: false,
	}
}

func (sw *MockStorageWatcher) AddTransformers(initializers []transformer.StorageTransformerInitializer) {
	fakeTransformer := &MockStorageTransformer{}
	sw.Transformers = []transformer.StorageTransformer{fakeTransformer}
	sw.AddTransformersWasCalled = true
}

func (sw *MockStorageWatcher) Execute(rows chan utils.StorageDiffRow, errs chan error) {
	sw.ExecuteWasCalled = true
}

func (sw *MockStorageWatcher) WatchEthStorage(wg *sync.WaitGroup) {
	defer wg.Done()
	sw.WatchEthStorageWasCalled = true
}
