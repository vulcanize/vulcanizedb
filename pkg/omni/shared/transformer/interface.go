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

package transformer

// Used to extract any/all events and a subset of method (state variable)
// data for any contract and persists it to custom postgres tables in vDB
type Transformer interface {
	SetEvents(contractAddr string, filterSet []string)
	SetEventArgs(contractAddr string, filterSet []string)
	SetMethods(contractAddr string, filterSet []string)
	SetMethodArgs(contractAddr string, filterSet []string)
	SetStartingBlock(contractAddr string, start int64)
	SetCreateAddrList(contractAddr string, on bool)
	SetCreateHashList(contractAddr string, on bool)
	SetPiping(contractAddr string, on bool)
	Init() error
	Execute() error
}
