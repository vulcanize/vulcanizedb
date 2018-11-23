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

package repository_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

var _ = Describe("Repository", func() {
	var db *postgres.DB
	var dataStore repository.MethodRepository
	var con *contract.Contract
	var err error
	var mockResult types.Result

	BeforeEach(func() {
		con = test_helpers.SetupTusdContract([]string{}, []string{"balanceOf"})
		method := con.Methods["balanceOf"]
		mockResult = types.Result{
			Method: method,
			PgType: method.Return[0].PgType,
			Inputs: make([]interface{}, 1),
			Output: new(interface{}),
			Block:  6707323,
		}
		mockResult.Inputs[0] = "0xfE9e8709d3215310075d67E3ed32A380CCf451C8"
		mockResult.Output = "66386309548896882859581786"
		db, _ = test_helpers.SetupDBandBC()
		dataStore = repository.NewMethodRepository(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("CreateContractSchema", func() {
		It("Creates schema if it doesn't exist", func() {
			created, err := dataStore.CreateContractSchema(constants.TusdContractAddress)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(true))

			created, err = dataStore.CreateContractSchema(constants.TusdContractAddress)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(false))
		})
	})

	Describe("CreateMethodTable", func() {
		It("Creates table if it doesn't exist", func() {
			created, err := dataStore.CreateContractSchema(constants.TusdContractAddress)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(true))

			created, err = dataStore.CreateMethodTable(constants.TusdContractAddress, mockResult)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(true))

			created, err = dataStore.CreateMethodTable(constants.TusdContractAddress, mockResult)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(Equal(false))
		})
	})

	Describe("PersistResult", func() {
		It("Persists result from method polling in custom pg table", func() {
			err = dataStore.PersistResult(mockResult, con.Address, con.Name)
			Expect(err).ToNot(HaveOccurred())

			scanStruct := test_helpers.BalanceOf{}

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM c%s.balanceof_method", constants.TusdContractAddress)).StructScan(&scanStruct)
			expectedLog := test_helpers.BalanceOf{
				Id:        1,
				TokenName: "TrueUSD",
				Block:     6707323,
				Address:   "0xfE9e8709d3215310075d67E3ed32A380CCf451C8",
				Balance:   "66386309548896882859581786",
			}
			Expect(scanStruct).To(Equal(expectedLog))
		})

		It("Fails with empty result", func() {
			err = dataStore.PersistResult(types.Result{}, con.Address, con.Name)
			Expect(err).To(HaveOccurred())
		})
	})
})
