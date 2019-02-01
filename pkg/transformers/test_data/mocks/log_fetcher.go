// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockLogFetcher struct {
	FetchedContractAddresses [][]common.Address
	FetchedTopics            [][]common.Hash
	FetchedBlocks            []int64
	fetcherError             error
	FetchedLogs              []types.Log
	SetBcCalled              bool
	FetchLogsCalled          bool
}

func (mlf *MockLogFetcher) FetchLogs(contractAddresses []common.Address, topics []common.Hash, header core.Header) ([]types.Log, error) {
	mlf.FetchedContractAddresses = append(mlf.FetchedContractAddresses, contractAddresses)
	mlf.FetchedTopics = [][]common.Hash{topics}
	mlf.FetchedBlocks = append(mlf.FetchedBlocks, header.BlockNumber)
	mlf.FetchLogsCalled = true

	return mlf.FetchedLogs, mlf.fetcherError
}

func (mlf *MockLogFetcher) SetBC(bc core.BlockChain) {
	mlf.SetBcCalled = true
}

func (mlf *MockLogFetcher) SetFetcherError(err error) {
	mlf.fetcherError = err
}

func (mlf *MockLogFetcher) SetFetchedLogs(logs []types.Log) {
	mlf.FetchedLogs = logs
}
