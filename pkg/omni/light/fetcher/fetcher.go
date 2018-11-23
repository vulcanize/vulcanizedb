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

type LogFetcher interface {
	FetchLogs(contractAddresses []string, topics [][]common.Hash, header core.Header) ([]types.Log, error)
}

type SettableLogFetcher interface {
	LogFetcher
	SetBC(bc core.BlockChain)
}

type Fetcher struct {
	blockChain core.BlockChain
}

func (fetcher *Fetcher) SetBC(bc core.BlockChain) {
	fetcher.blockChain = bc
}

func NewFetcher(blockchain core.BlockChain) Fetcher {
	return Fetcher{
		blockChain: blockchain,
	}
}

func (fetcher Fetcher) FetchLogs(contractAddresses []string, topics [][]common.Hash, header core.Header) ([]types.Log, error) {
	addresses := hexStringsToAddresses(contractAddresses)
	blockHash := common.HexToHash(header.Hash)
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
		Addresses: addresses,
		Topics:    topics,
	}

	return fetcher.blockChain.GetEthLogsWithCustomQuery(query)
}

func hexStringsToAddresses(hexStrings []string) []common.Address {
	var addresses []common.Address
	for _, hexString := range hexStrings {
		address := common.HexToAddress(hexString)
		addresses = append(addresses, address)
	}

	return addresses
}
