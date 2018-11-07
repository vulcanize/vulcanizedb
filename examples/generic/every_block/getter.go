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
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

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
	generic.Fetcher // Underlying Fetcher
}

// Initializes and returns a Getter with the given blockchain
func NewGetter(blockChain core.BlockChain) GenericGetter {
	return GenericGetter{
		Fetcher: generic.Fetcher{
			BlockChain: blockChain,
		},
	}
}

// Public getter methods for calling contract methods
func (g GenericGetter) GetOwner(contractAbi, contractAddress string, blockNumber int64) (common.Address, error) {
	return g.Fetcher.FetchAddress("owner", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetStoppedStatus(contractAbi, contractAddress string, blockNumber int64) (bool, error) {
	return g.Fetcher.FetchBool("stopped", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetStringName(contractAbi, contractAddress string, blockNumber int64) (string, error) {
	return g.Fetcher.FetchString("name", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetHashName(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error) {
	return g.Fetcher.FetchHash("name", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetStringSymbol(contractAbi, contractAddress string, blockNumber int64) (string, error) {
	return g.Fetcher.FetchString("symbol", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetHashSymbol(contractAbi, contractAddress string, blockNumber int64) (common.Hash, error) {
	return g.Fetcher.FetchHash("symbol", contractAbi, contractAddress, blockNumber, nil)
}

func (g GenericGetter) GetDecimals(contractAbi, contractAddress string, blockNumber int64) (big.Int, error) {
	return g.Fetcher.FetchBigInt("decimals", contractAbi, contractAddress, blockNumber, nil)
}

// Method to retrieve the Getter's blockchain
func (g GenericGetter) GetBlockChain() core.BlockChain {
	return g.Fetcher.BlockChain
}
