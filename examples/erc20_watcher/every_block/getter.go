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
	generic.Fetcher
}

// Initializes and returns a Getter with the given blockchain
func NewGetter(blockChain core.BlockChain) ERC20Getter {
	return ERC20Getter{
		Fetcher: generic.Fetcher{
			BlockChain: blockChain,
		},
	}
}

// Public getter methods for calling contract methods
func (g ERC20Getter) GetTotalSupply(contractAbi, contractAddress string, blockNumber int64) (big.Int, error) {
	return g.Fetcher.FetchBigInt("totalSupply", contractAbi, contractAddress, blockNumber, nil)
}

func (g ERC20Getter) GetBalance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	return g.Fetcher.FetchBigInt("balanceOf", contractAbi, contractAddress, blockNumber, methodArgs)
}

func (g ERC20Getter) GetAllowance(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (big.Int, error) {
	return g.Fetcher.FetchBigInt("allowance", contractAbi, contractAddress, blockNumber, methodArgs)
}

// Method to retrieve the Getter's blockchain
func (g ERC20Getter) GetBlockChain() core.BlockChain {
	return g.Fetcher.BlockChain
}
