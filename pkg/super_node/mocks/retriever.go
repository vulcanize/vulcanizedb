package mocks

import (
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// MockCIDRetriever is a mock CID retriever for use in tests
type MockCIDRetriever struct {
	GapsToRetrieve              [][2]int64
	GapsToRetrieveErr           error
	CalledTimes                 int
	FirstBlockNumberToReturn    int64
	RetrieveFirstBlockNumberErr error
}

// RetrieveCIDs mock method
func (*MockCIDRetriever) RetrieveCIDs(streamFilters config.Subscription, blockNumber int64) (*ipfs.CIDWrapper, error) {
	panic("implement me")
}

// RetrieveLastBlockNumber mock method
func (*MockCIDRetriever) RetrieveLastBlockNumber() (int64, error) {
	panic("implement me")
}

// RetrieveFirstBlockNumber mock method
func (mcr *MockCIDRetriever) RetrieveFirstBlockNumber() (int64, error) {
	return mcr.FirstBlockNumberToReturn, mcr.RetrieveFirstBlockNumberErr
}

// RetrieveGapsInData mock method
func (mcr *MockCIDRetriever) RetrieveGapsInData() ([][2]int64, error) {
	mcr.CalledTimes++
	return mcr.GapsToRetrieve, mcr.GapsToRetrieveErr
}

// SetGapsToRetrieve mock method
func (mcr *MockCIDRetriever) SetGapsToRetrieve(gaps [][2]int64) {
	if mcr.GapsToRetrieve == nil {
		mcr.GapsToRetrieve = make([][2]int64, 0)
	}
	mcr.GapsToRetrieve = append(mcr.GapsToRetrieve, gaps...)
}
