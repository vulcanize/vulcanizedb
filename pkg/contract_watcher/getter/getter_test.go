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

package getter_test

import (
	"math/rand"

	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/getter"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Interface Getter", func() {
	Describe("GetAbi", func() {
		It("fetches the contract's data from the blockchain", func() {
			blockChain := fakes.MockBlockChain{}
			g := getter.NewInterfaceGetter(&blockChain)
			testAddress := fakes.FakeAddress.Hex()
			testBlockNumber := rand.Int63()
			g.GetABI(testAddress, testBlockNumber)

			methodArgs := make([]interface{}, 1)
			methodArgs[0] = constants.MetaSig.Bytes()
			result := new(bool)

			blockChain.AssertFetchContractDataCalledWith(constants.SupportsInterfaceABI, testAddress,
				"supportsInterface", methodArgs, &result, testBlockNumber)
		})
	})
})
