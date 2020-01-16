package mocks

import (
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// MockCIDRetriever is a mock CID retriever for use in tests
type MockCIDRetriever struct {
	GapsToRetrieve              [][2]uint64
	GapsToRetrieveErr           error
	CalledTimes                 int
	FirstBlockNumberToReturn    int64
	RetrieveFirstBlockNumberErr error
}

// RetrieveCIDs mock method
func (*MockCIDRetriever) RetrieveCIDs(streamFilters config.Subscription, blockNumber int64) (*ipfs.CIDWrapper, error) {
	panic("implement me")
}

// RetrieveHeaderCIDs mock method
func (*MockCIDRetriever) RetrieveHeaderCIDs(tx *sqlx.Tx, blockNumber int64) ([]string, error) {
	panic("implement me")

}

// RetrieveUncleCIDs mock method
func (*MockCIDRetriever) RetrieveUncleCIDs(tx *sqlx.Tx, blockNumber int64) ([]string, error) {
	panic("implement me")

}

// RetrieveTrxCIDs mock method
func (*MockCIDRetriever) RetrieveTrxCIDs(tx *sqlx.Tx, txFilter config.TrxFilter, blockNumber int64) ([]string, []int64, error) {
	panic("implement me")

}

// RetrieveRctCIDs mock method
func (*MockCIDRetriever) RetrieveRctCIDs(tx *sqlx.Tx, rctFilter config.ReceiptFilter, blockNumber int64, trxIds []int64) ([]string, error) {
	panic("implement me")

}

// RetrieveStateCIDs mock method
func (*MockCIDRetriever) RetrieveStateCIDs(tx *sqlx.Tx, stateFilter config.StateFilter, blockNumber int64) ([]ipfs.StateNodeCID, error) {
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
func (mcr *MockCIDRetriever) RetrieveGapsInData() ([][2]uint64, error) {
	mcr.CalledTimes++
	return mcr.GapsToRetrieve, mcr.GapsToRetrieveErr
}

// SetGapsToRetrieve mock method
func (mcr *MockCIDRetriever) SetGapsToRetrieve(gaps [][2]uint64) {
	if mcr.GapsToRetrieve == nil {
		mcr.GapsToRetrieve = make([][2]uint64, 0)
	}
	mcr.GapsToRetrieve = append(mcr.GapsToRetrieve, gaps...)
}

func (mcr *MockCIDRetriever) Database() *postgres.DB {
	panic("implement me")
}
