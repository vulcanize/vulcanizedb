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
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type LogFetcher interface {
	FetchLogs(contractAddresses []string, topics [][]common.Hash, blockNumber int64) ([]types.Log, error)
}

type Fetcher struct {
	blockChain core.BlockChain
}

func NewFetcher(blockchain core.BlockChain) Fetcher {
	return Fetcher{
		blockChain: blockchain,
	}
}

func (fetcher Fetcher) FetchLogs(contractAddresses []string, topics [][]common.Hash, blockNumber int64) ([]types.Log, error) {
	block := big.NewInt(blockNumber)
	addresses := hexStringsToAddresses(contractAddresses)
	query := ethereum.FilterQuery{
		FromBlock: block,
		ToBlock:   block,
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