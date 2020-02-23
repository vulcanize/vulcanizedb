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

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
)

type EthHeaderDagPutter struct {
	adder *ipfs.IPFS
}

func NewEthBlockHeaderDagPutter(adder *ipfs.IPFS) *EthHeaderDagPutter {
	return &EthHeaderDagPutter{adder: adder}
}

func (bhdp *EthHeaderDagPutter) DagPut(raw interface{}) ([]string, error) {
	header, ok := raw.(*types.Header)
	if !ok {
		return nil, fmt.Errorf("EthHeaderDagPutter expected input type %T got %T", &types.Header{}, raw)
	}
	node, err := ipld.NewEthHeader(header)
	if err != nil {
		return nil, err
	}
	if err := bhdp.adder.Add(node); err != nil {
		return nil, err
	}
	return []string{node.Cid().String()}, nil
}
