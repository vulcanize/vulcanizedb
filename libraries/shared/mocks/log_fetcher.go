// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	ContractAddresses []common.Address
	FetchCalled       bool
	MissingHeader     core.Header
	ReturnError       error
	ReturnLogs        []types.Log
	Topics            []common.Hash
}

func (fetcher *MockLogFetcher) FetchLogs(contractAddresses []common.Address, topics []common.Hash, missingHeader core.Header) ([]types.Log, error) {
	fetcher.FetchCalled = true
	fetcher.ContractAddresses = contractAddresses
	fetcher.Topics = topics
	fetcher.MissingHeader = missingHeader
	return fetcher.ReturnLogs, fetcher.ReturnError
}
