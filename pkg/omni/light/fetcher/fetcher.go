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

package fetcher

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Fetcher interface {
	FetchLogs(contractAddresses []string, topics []common.Hash, missingHeader core.Header) ([]types.Log, error)
}

type fetcher struct {
	blockChain core.BlockChain
}

func NewFetcher(blockchain core.BlockChain) *fetcher {
	return &fetcher{
		blockChain: blockchain,
	}
}

// Checks all topic0s, on all addresses, fetching matching logs for the given header
func (fetcher *fetcher) FetchLogs(contractAddresses []string, topic0s []common.Hash, header core.Header) ([]types.Log, error) {
	addresses := hexStringsToAddresses(contractAddresses)
	blockHash := common.HexToHash(header.Hash)
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: addresses,
		// Search for _any_ of the topics in topic0 position; see docs on `FilterQuery`
		Topics: [][]common.Hash{topic0s},
	}

	logs, err := fetcher.blockChain.GetEthLogsWithCustomQuery(query)
	if err != nil {
		// TODO review aggregate fetching error handling
		return []types.Log{}, err
	}

	return logs, nil
}

func hexStringsToAddresses(hexStrings []string) []common.Address {
	var addresses []common.Address
	for _, hexString := range hexStrings {
		address := common.HexToAddress(hexString)
		addresses = append(addresses, address)
	}

	return addresses
}
