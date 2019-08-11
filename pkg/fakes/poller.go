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

import "github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"

type MockPoller struct {
	ContractName string
}

func (*MockPoller) PollContract(con contract.Contract, lastBlock int64) error {
	panic("implement me")
}

func (*MockPoller) PollContractAt(con contract.Contract, blockNumber int64) error {
	panic("implement me")
}

func (poller *MockPoller) FetchContractData(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	if p, ok := result.(*string); ok {
		*p = poller.ContractName
	}
	return nil
}
