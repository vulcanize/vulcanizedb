// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy og the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS Og ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package every_block

import (
	"math/big"

	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// Getter serves as a higher level data fetcher that invokes its underlying Fetcher methods for a given contract method

// Interface definition for a Getter
type ERC20GetterInterface interface {
	GetTotalSupply(contractAbi, contractAddress string, blockNumber int64) (big.Int, error)
	GetBalance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error)
	GetAllowance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error)
	GetBlockChain() core.BlockChain
}

// Getter struct
type ERC20Getter struct {
	fetcher generic.Fetcher
}

// Initializes and returns a Getter with the given blockchain
func NewGetter(blockChain core.BlockChain) ERC20Getter {
	return ERC20Getter{
		fetcher: generic.Fetcher{
			BlockChain: blockChain,
		},
	}
}

// Public getter methods for calling contract methods
func (g ERC20Getter) GetTotalSupply(contractAbi, contractAddress string, blockNumber int64) (big.Int, error) {
	return g.fetcher.FetchBigInt("totalSupply", contractAbi, contractAddress, blockNumber, nil)
}

func (g ERC20Getter) GetBalance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	return g.fetcher.FetchBigInt("balanceOf", contractAbi, contractAddress, blockNumber, methodArgs)
}

func (g ERC20Getter) GetAllowance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	return g.fetcher.FetchBigInt("allowance", contractAbi, contractAddress, blockNumber, methodArgs)
}

// Method to retrieve the Getter's blockchain
func (g ERC20Getter) GetBlockChain() core.BlockChain {
	return g.fetcher.BlockChain
}
