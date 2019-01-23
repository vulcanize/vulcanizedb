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

package getter

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/fetcher"
)

type InterfaceGetter interface {
	GetABI(resolverAddr string, blockNumber int64) string
	GetBlockChain() core.BlockChain
}

type interfaceGetter struct {
	fetcher.Fetcher
}

func NewInterfaceGetter(blockChain core.BlockChain) *interfaceGetter {
	return &interfaceGetter{
		Fetcher: fetcher.Fetcher{
			BlockChain: blockChain,
		},
	}
}

// Used to construct a custom ABI based on the results from calling supportsInterface
func (g *interfaceGetter) GetABI(resolverAddr string, blockNumber int64) string {
	a := constants.SupportsInterfaceABI
	args := make([]interface{}, 1)
	args[0] = constants.MetaSig.Bytes()
	supports, err := g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err != nil || !supports {
		return ""
	}
	abiStr := `[`
	args[0] = constants.AddrChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.AddrChangeInterface + ","
	}
	args[0] = constants.NameChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.NameChangeInterface + ","
	}
	args[0] = constants.ContentChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.ContentChangeInterface + ","
	}
	args[0] = constants.AbiChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.AbiChangeInterface + ","
	}
	args[0] = constants.PubkeyChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.PubkeyChangeInterface + ","
	}
	args[0] = constants.ContentHashChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.ContenthashChangeInterface + ","
	}
	args[0] = constants.MultihashChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.MultihashChangeInterface + ","
	}
	args[0] = constants.TextChangeSig.Bytes()
	supports, err = g.getSupportsInterface(a, resolverAddr, blockNumber, args)
	if err == nil && supports {
		abiStr += constants.TextChangeInterface + ","
	}
	abiStr = abiStr[:len(abiStr)-1] + `]`

	return abiStr
}

// Use this method to check whether or not a contract supports a given method/event interface
func (g *interfaceGetter) getSupportsInterface(contractAbi, contractAddress string, blockNumber int64, methodArgs []interface{}) (bool, error) {
	return g.Fetcher.FetchBool("supportsInterface", contractAbi, contractAddress, blockNumber, methodArgs)
}

// Method to retrieve the Getter's blockchain
func (g *interfaceGetter) GetBlockChain() core.BlockChain {
	return g.Fetcher.BlockChain
}
