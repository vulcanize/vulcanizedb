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

package flip_kick

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type LogFetcher interface {
	FetchLogs(contractAddress string, topics [][]common.Hash, blockNumber int64) ([]types.Log, error)
}

type Fetcher struct {
	Blockchain core.BlockChain
}

func NewFetcher(blockchain core.BlockChain) Fetcher {
	return Fetcher{
		Blockchain: blockchain,
	}
}

func (f Fetcher) FetchLogs(contractAddress string, topicZeros [][]common.Hash, blockNumber int64) ([]types.Log, error) {
	block := big.NewInt(blockNumber)
	address := common.HexToAddress(contractAddress)
	query := ethereum.FilterQuery{
		FromBlock: block,
		ToBlock:   block,
		Addresses: []common.Address{address},
		Topics:    topicZeros,
	}

	return f.Blockchain.GetEthLogsWithCustomQuery(query)
}
