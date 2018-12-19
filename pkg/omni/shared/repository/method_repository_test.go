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
	"strings"

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
	var method types.Method

	BeforeEach(func() {
		con = test_helpers.SetupTusdContract([]string{}, []string{"balanceOf"})
		Expect(len(con.Methods)).To(Equal(1))
		method = con.Methods[0]
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
		dataStore = repository.NewMethodRepository(db, types.FullSync)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("Full Sync Mode", func() {
		BeforeEach(func() {
			dataStore = repository.NewMethodRepository(db, types.FullSync)
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

			It("Caches schema it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				_, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(false))

				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("CreateMethodTable", func() {
			It("Creates table if it doesn't exist", func() {
				created, err := dataStore.CreateContractSchema(constants.TusdContractAddress)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateMethodTable(constants.TusdContractAddress, method)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateMethodTable(constants.TusdContractAddress, method)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(false))
			})

			It("Caches table it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				tableID := fmt.Sprintf("%s_%s.%s_method", types.FullSync, strings.ToLower(con.Address), strings.ToLower(method.Name))
				_, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(false))

				created, err = dataStore.CreateMethodTable(con.Address, method)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("PersistResult", func() {
			It("Persists result from method polling in custom pg table", func() {
				err = dataStore.PersistResults([]types.Result{mockResult}, method, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.BalanceOf{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method", constants.TusdContractAddress)).StructScan(&scanStruct)
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
				err = dataStore.PersistResults([]types.Result{}, method, con.Address, con.Name)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Light Sync Mode", func() {
		BeforeEach(func() {
			dataStore = repository.NewMethodRepository(db, types.LightSync)
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

			It("Caches schema it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				_, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(false))

				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckSchemaCache(con.Address)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("CreateMethodTable", func() {
			It("Creates table if it doesn't exist", func() {
				created, err := dataStore.CreateContractSchema(constants.TusdContractAddress)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateMethodTable(constants.TusdContractAddress, method)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				created, err = dataStore.CreateMethodTable(constants.TusdContractAddress, method)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(false))
			})

			It("Caches table it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
				created, err := dataStore.CreateContractSchema(con.Address)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				tableID := fmt.Sprintf("%s_%s.%s_method", types.LightSync, strings.ToLower(con.Address), strings.ToLower(method.Name))
				_, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(false))

				created, err = dataStore.CreateMethodTable(con.Address, method)
				Expect(err).ToNot(HaveOccurred())
				Expect(created).To(Equal(true))

				v, ok := dataStore.CheckTableCache(tableID)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			})
		})

		Describe("PersistResult", func() {
			It("Persists result from method polling in custom pg table for light sync mode vDB", func() {
				err = dataStore.PersistResults([]types.Result{mockResult}, method, con.Address, con.Name)
				Expect(err).ToNot(HaveOccurred())

				scanStruct := test_helpers.BalanceOf{}

				err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method", constants.TusdContractAddress)).StructScan(&scanStruct)
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
				err = dataStore.PersistResults([]types.Result{}, method, con.Address, con.Name)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
