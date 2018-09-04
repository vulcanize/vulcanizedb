// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"math/big"
)

type IPriceFeedFetcher interface {
	FetchLogValues(blockNumber int64) ([]types.Log, error)
}

type PriceFeedFetcher struct {
	blockChain        core.BlockChain
	contractAddresses []string
}

func NewPriceFeedFetcher(blockChain core.BlockChain, contractAddresses []string) PriceFeedFetcher {
	return PriceFeedFetcher{
		blockChain:        blockChain,
		contractAddresses: contractAddresses,
	}
}

func (fetcher PriceFeedFetcher) FetchLogValues(blockNumber int64) ([]types.Log, error) {
	var addresses []common.Address
	for _, addr := range fetcher.contractAddresses {
		addresses = append(addresses, common.HexToAddress(addr))
	}
	n := big.NewInt(blockNumber)
	query := ethereum.FilterQuery{
		FromBlock: n,
		ToBlock:   n,
		Addresses: addresses,
		Topics:    [][]common.Hash{{common.HexToHash(shared.LogValueSignature)}},
	}
	return fetcher.blockChain.GetEthLogsWithCustomQuery(query)
}
