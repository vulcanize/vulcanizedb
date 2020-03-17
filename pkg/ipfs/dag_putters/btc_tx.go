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

	node "github.com/ipfs/go-ipld-format"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
)

type BtcTxDagPutter struct {
	adder *ipfs.IPFS
}

func NewBtcTxDagPutter(adder *ipfs.IPFS) *BtcTxDagPutter {
	return &BtcTxDagPutter{adder: adder}
}

func (etdp *BtcTxDagPutter) DagPut(n node.Node) (string, error) {
	transaction, ok := n.(*ipld.BtcTx)
	if !ok {
		return "", fmt.Errorf("BtcTxDagPutter expected input type %T got %T", &ipld.BtcTx{}, n)
	}
	if err := etdp.adder.Add(transaction); err != nil {
		return "", err
	}
	return transaction.Cid().String(), nil
}
