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

package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Getter struct {
	Fetcher Fetcher
}

func NewGetter(blockChain core.BlockChain) Getter {
	return Getter{
		Fetcher: Fetcher{
			BlockChain: blockChain,
		},
	}
}

func (g *Getter) GetTotalSupply(contractAbi, contractAddress string, blockNumber int64) (big.Int, error) {
	return g.Fetcher.FetchBigInt("totalSupply", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetBalance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	return g.Fetcher.FetchBigInt("balanceOf", contractAbi, contractAddress, blockNumber, methodArgs)
}

func (g *Getter) GetAllowance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	return g.Fetcher.FetchBigInt("allowance", contractAbi, contractAddress, blockNumber, methodArgs)
}

func (g *Getter) GetOwner(contractAbi, contractAddress string, blockNumber int64) (common.Address, error) {
	return g.Fetcher.FetchAddress("owner", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetStoppedStatus(contractAbi, contractAddress string, blockNumber int64) (bool, error) {
	return g.Fetcher.FetchBool("stopped", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetStringName(contractAbi, contractAddress string, blockNumber int64) (string, error) {
	return g.Fetcher.FetchString("name", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetHashName(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error) {
	return g.Fetcher.FetchHash("name", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetStringSymbol(contractAbi, contractAddress string, blockNumber int64) (string, error) {
	return g.Fetcher.FetchString("symbol", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetHashSymbol(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error) {
	return g.Fetcher.FetchHash("symbol", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetDecimals(contractAbi, contractAddress string, blockNumber int64) (big.Int, error) {
	return g.Fetcher.FetchBigInt("decimals", contractAbi, contractAddress, blockNumber, nil)
}

func (g *Getter) GetBlockChain() core.BlockChain {
	return g.Fetcher.BlockChain
}
