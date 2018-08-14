package price_feeds

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockPriceFeedFetcher struct {
	passedBlockNumbers []int64
	returnErr          error
	returnLogs         []types.Log
}

func (fetcher *MockPriceFeedFetcher) SetReturnErr(err error) {
	fetcher.returnErr = err
}

func (fetcher *MockPriceFeedFetcher) SetReturnLogs(logs []types.Log) {
	fetcher.returnLogs = logs
}

func (fetcher *MockPriceFeedFetcher) FetchLogValues(blockNumber int64) ([]types.Log, error) {
	fetcher.passedBlockNumbers = append(fetcher.passedBlockNumbers, blockNumber)
	return fetcher.returnLogs, fetcher.returnErr
}

func (fetcher *MockPriceFeedFetcher) AssertFetchLogValuesCalledWith(blockNumbers []int64) {
	Expect(fetcher.passedBlockNumbers).To(Equal(blockNumbers))
}
