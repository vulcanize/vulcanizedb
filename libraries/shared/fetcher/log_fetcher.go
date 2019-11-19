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

package fetcher

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/makerdao/vulcanizedb/pkg/core"
)

type ILogFetcher interface {
	FetchLogs(contractAddresses []common.Address, topics []common.Hash, missingHeader core.Header) ([]types.Log, error)
	// TODO Extend FetchLogs for doing several blocks at a time
}

type LogFetcher struct {
	blockChain core.BlockChain
}

func NewLogFetcher(blockchain core.BlockChain) *LogFetcher {
	return &LogFetcher{
		blockChain: blockchain,
	}
}

// Checks all topic0s, on all addresses, fetching matching logs for the given header
func (logFetcher LogFetcher) FetchLogs(addresses []common.Address, topic0s []common.Hash, header core.Header) ([]types.Log, error) {
	blockHash := common.HexToHash(header.Hash)
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: addresses,
		// Search for _any_ of the topics in topic0 position; see docs on `FilterQuery`
		Topics: [][]common.Hash{topic0s},
	}

	logs, err := logFetcher.blockChain.GetEthLogsWithCustomQuery(query)
	if err != nil {
		// TODO review aggregate fetching error handling
		return []types.Log{}, err
	}

	return logs, nil
}
