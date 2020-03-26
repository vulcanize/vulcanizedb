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
	"strings"

	node "github.com/ipfs/go-ipld-format"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
)

type EthReceiptDagPutter struct {
	adder *ipfs.IPFS
}

func NewEthReceiptDagPutter(adder *ipfs.IPFS) *EthReceiptDagPutter {
	return &EthReceiptDagPutter{adder: adder}
}

func (erdp *EthReceiptDagPutter) DagPut(n node.Node) (string, error) {
	receipt, ok := n.(*ipld.EthReceipt)
	if !ok {
		return "", fmt.Errorf("EthReceiptDagPutter expected input type %T got type %T", &ipld.EthReceipt{}, n)
	}
	if err := erdp.adder.Add(receipt); err != nil && !strings.Contains(err.Error(), duplicateKeyErrorString) {
		return "", err
	}
	return receipt.Cid().String(), nil
}
