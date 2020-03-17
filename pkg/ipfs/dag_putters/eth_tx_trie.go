// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package dag_putters

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
)

type EthTxTrieDagPutter struct {
	adder *ipfs.IPFS
}

func NewEthTxTrieDagPutter(adder *ipfs.IPFS) *EthTxTrieDagPutter {
	return &EthTxTrieDagPutter{adder: adder}
}

func (etdp *EthTxTrieDagPutter) DagPut(raw interface{}) ([]string, error) {
	txTrieNodes, ok := raw.([]*ipld.EthTxTrie)
	if !ok {
		return nil, fmt.Errorf("EthTxTrieDagPutter expected input type %T got %T", []*ipld.EthTxTrie{}, raw)
	}
	cids := make([]string, len(txTrieNodes))
	for i, txNode := range txTrieNodes {
		if err := etdp.adder.Add(txNode); err != nil {
			return nil, err
		}
		cids[i] = txNode.Cid().String()
	}
	return cids, nil
}
