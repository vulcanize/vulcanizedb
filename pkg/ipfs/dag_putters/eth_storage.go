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

type EthStorageDagPutter struct {
	adder *ipfs.IPFS
}

func NewEthStorageDagPutter(adder *ipfs.IPFS) *EthStorageDagPutter {
	return &EthStorageDagPutter{adder: adder}
}

func (erdp *EthStorageDagPutter) DagPut(raw interface{}) ([]string, error) {
	storageNode, ok := raw.(*ipld.EthStorageTrie)
	if !ok {
		return nil, fmt.Errorf("EthStorageDagPutter expected input type %T got %T", &ipld.EthStorageTrie{}, raw)
	}
	if err := erdp.adder.Add(storageNode); err != nil {
		return nil, err
	}
	return []string{storageNode.Cid().String()}, nil
}
