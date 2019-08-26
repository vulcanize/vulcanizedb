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

import "errors"

// DagPutter is a mock for testing the ipfs publisher
type DagPutter struct {
	CIDsToReturn []string
	ErrToReturn  error
}

// DagPut returns the pre-loaded CIDs or error
func (dp *DagPutter) DagPut(raw interface{}) ([]string, error) {
	return dp.CIDsToReturn, dp.ErrToReturn
}

// IncrementingDagPutter is a mock for testing the ipfs publisher
type IncrementingDagPutter struct {
	CIDsToReturn []string
	iterator     int
	ErrToReturn  error
}

// DagPut returns the pre-loaded CIDs or error
func (dp *IncrementingDagPutter) DagPut(raw interface{}) ([]string, error) {
	if len(dp.CIDsToReturn) >= dp.iterator+1 {
		cid := dp.CIDsToReturn[dp.iterator]
		dp.iterator++
		return []string{cid}, dp.ErrToReturn
	}
	return nil, errors.New("dag putter iterator is out of range")
}
