package mocks

import "github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"

type MockStorageFetcher struct {
	RowsToReturn []utils.StorageDiffRow
	ErrsToReturn []error
}

func NewMockStorageFetcher() *MockStorageFetcher {
	return &MockStorageFetcher{}
}

func (fetcher *MockStorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiffRow, errs chan<- error) {
	for _, err := range fetcher.ErrsToReturn {
		errs <- err
	}
	for _, row := range fetcher.RowsToReturn {
		out <- row
	}
	close(out)
	close(errs)
}
