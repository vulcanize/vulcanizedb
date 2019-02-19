package mocks

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

type MockStorageQueue struct {
	AddCalled bool
	AddError  error
}

func (queue *MockStorageQueue) Add(row shared.StorageDiffRow) error {
	queue.AddCalled = true
	return queue.AddError
}
