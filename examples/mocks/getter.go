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
