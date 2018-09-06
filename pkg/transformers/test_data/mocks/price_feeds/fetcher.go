// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
