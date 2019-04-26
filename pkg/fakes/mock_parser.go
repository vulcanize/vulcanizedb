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

package fakes

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/types"
)

type MockParser struct {
	AbiToReturn string
	EventName   string
	Event       types.Event
}

func (*MockParser) Parse(contractAddr string) error {
	return nil
}

func (m *MockParser) ParseAbiStr(abiStr string) error {
	m.AbiToReturn = abiStr
	return nil
}

func (parser *MockParser) Abi() string {
	return parser.AbiToReturn
}

func (*MockParser) ParsedAbi() abi.ABI {
	return abi.ABI{}
}

func (*MockParser) GetMethods(wanted []string) []types.Method {
	panic("implement me")
}

func (*MockParser) GetSelectMethods(wanted []string) []types.Method {
	return []types.Method{}
}

func (parser *MockParser) GetEvents(wanted []string) map[string]types.Event {
	return map[string]types.Event{parser.EventName: parser.Event}
}
