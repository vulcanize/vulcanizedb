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

package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MockLogFetcher struct {
	FetchedContractAddress string
	FetchedTopics          [][]common.Hash
	FetchedBlocks          []int64
	fetcherError           error
	FetchedLogs            []types.Log
}

func (mlf *MockLogFetcher) FetchLogs(contractAddress string, topics [][]common.Hash, blockNumber int64) ([]types.Log, error) {
	mlf.FetchedContractAddress = contractAddress
	mlf.FetchedTopics = topics
	mlf.FetchedBlocks = append(mlf.FetchedBlocks, blockNumber)

	return mlf.FetchedLogs, mlf.fetcherError
}

func (mlf *MockLogFetcher) SetFetcherError(err error) {
	mlf.fetcherError = err
}

func (mlf *MockLogFetcher) SetFetchedLogs(logs []types.Log) {
	mlf.FetchedLogs = logs
}
