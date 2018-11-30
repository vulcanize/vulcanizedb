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

package shared

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// TODO Check if Fetcher can be simplified with aggregate logic

type LogFetcher interface {
	FetchLogs(contractAddresses []common.Address, topics []common.Hash, missingHeader core.Header) ([]types.Log, error)
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

// Checks all topic0s, on all addresses, fetching matching logs for the given header
func (fetcher Fetcher) FetchLogs(addresses []common.Address, topic0s []common.Hash, header core.Header) ([]types.Log, error) {
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
