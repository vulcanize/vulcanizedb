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

package every_block

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// Getter serves as a higher level data fetcher that invokes its underlying Fetcher methods for a given contract method

// Interface definition for a Getter
type GenericGetterInterface interface {
	GetOwner(contractAbi, contractAddress string, blockNumber int64) (common.Address, error)
	GetStoppedStatus(contractAbi, contractAddress string, blockNumber int64) (bool, error)
	GetStringName(contractAbi, contractAddress string, blockNumber int64) (string, error)
	GetHashName(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error)
	GetStringSymbol(contractAbi, contractAddress string, blockNumber int64) (string, error)
	GetHashSymbol(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error)
	GetDecimals(contractAbi, contractAddress string, blockNumber int64) (big.Int, error)
	GetBlockChain() core.BlockChain
}

// Getter struct
type GenericGetter struct {
	fetcher generic.Fetcher // Underlying Fetcher
}

// Initializes and returns a Getter with the given blockchain
func NewGetter(blockChain core.BlockChain) GenericGetter {
	return GenericGetter{
		fetcher: generic.Fetcher{
			BlockChain: blockChain,
		},
	}
}

// Public getter methods for calling contract methods
func (g GenericGetter) GetOwner(contractAbi, contractAddress string, blockNumber int64) (common.Address, error) {
	return g.fetcher.FetchAddress("owner", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetStoppedStatus(contractAbi, contractAddress string, blockNumber int64) (bool, error) {
	return g.fetcher.FetchBool("stopped", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetStringName(contractAbi, contractAddress string, blockNumber int64) (string, error) {
	return g.fetcher.FetchString("name", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetHashName(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error) {
	return g.fetcher.FetchHash("name", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetStringSymbol(contractAbi, contractAddress string, blockNumber int64) (string, error) {
	return g.fetcher.FetchString("symbol", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetHashSymbol(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error) {
	return g.fetcher.FetchHash("symbol", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetDecimals(contractAbi, contractAddress string, blockNumber int64) (big.Int, error) {
	return g.fetcher.FetchBigInt("decimals", contractAbi, contractAddress, blockNumber, nil)
}

// Method to retrieve the Getter's blockchain
func (g GenericGetter) GetBlockChain() core.BlockChain {
	return g.fetcher.BlockChain
}
