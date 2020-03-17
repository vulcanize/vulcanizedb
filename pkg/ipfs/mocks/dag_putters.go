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

package mocks

import (
	"errors"

	node "github.com/ipfs/go-ipld-format"

	"github.com/ethereum/go-ethereum/common"
)

// DagPutter is a mock for testing the ipfs publisher
type DagPutter struct {
	PassedNode  node.Node
	ErrToReturn error
}

// DagPut returns the pre-loaded CIDs or error
func (dp *DagPutter) DagPut(n node.Node) (string, error) {
	dp.PassedNode = n
	return n.Cid().String(), dp.ErrToReturn
}

// MappedDagPutter is a mock for testing the ipfs publisher
type MappedDagPutter struct {
	CIDsToReturn map[common.Hash]string
	PassedNode   node.Node
	ErrToReturn  error
}

// DagPut returns the pre-loaded CIDs or error
func (mdp *MappedDagPutter) DagPut(n node.Node) (string, error) {
	if mdp.CIDsToReturn == nil {
		return "", errors.New("mapped dag putter needs to be initialized with a map of cids to return")
	}
	hash := common.BytesToHash(n.RawData())
	return mdp.CIDsToReturn[hash], nil
}
