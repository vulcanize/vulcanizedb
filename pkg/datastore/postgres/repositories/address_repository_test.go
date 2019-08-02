// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repositories_test

import (
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
	"strings"
)

var _ = Describe("address lookup", func() {
	var (
		db      *postgres.DB
		repo    repositories.AddressRepository
		address = fakes.FakeAddress.Hex()
	)
	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = repositories.AddressRepository{}
	})

	type dbAddress struct {
		Id      int
		Address string
	}

	Describe("GetOrCreateAddress", func() {
		It("creates an address record", func() {
			addressId, createErr := repo.GetOrCreateAddress(db, address)
			Expect(createErr).NotTo(HaveOccurred())

			var actualAddress dbAddress
			getErr := db.Get(&actualAddress, `SELECT id, address FROM public.addresses LIMIT 1`)
			Expect(getErr).NotTo(HaveOccurred())
			expectedAddress := dbAddress{Id: addressId, Address: address}
			Expect(actualAddress).To(Equal(expectedAddress))
		})

		It("returns the existing record id if the address already exists", func() {
			_, createErr := repo.GetOrCreateAddress(db, address)
			Expect(createErr).NotTo(HaveOccurred())

			_, getErr := repo.GetOrCreateAddress(db, address)
			Expect(getErr).NotTo(HaveOccurred())

			var addressCount int
			addressErr := db.Get(&addressCount, `SELECT count(*) FROM public.addresses`)
			Expect(addressErr).NotTo(HaveOccurred())
		})

		It("gets upper-cased addresses", func() {
			upperAddress := strings.ToUpper(address)
			upperAddressId, createErr := repo.GetOrCreateAddress(db, upperAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repo.GetOrCreateAddress(db, address)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})

		It("gets lower-cased addresses", func() {
			lowerAddress := strings.ToLower(address)
			upperAddressId, createErr := repo.GetOrCreateAddress(db, lowerAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repo.GetOrCreateAddress(db, address)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})
	})

	Describe("GetOrCreateAddressInTransaction", func() {
		var (
			tx    *sqlx.Tx
			txErr error
		)
		BeforeEach(func() {
			tx, txErr = db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())
		})

		It("creates an address record", func() {
			addressId, createErr := repo.GetOrCreateAddressInTransaction(tx, address)
			Expect(createErr).NotTo(HaveOccurred())
			tx.Commit()

			var actualAddress dbAddress
			getErr := db.Get(&actualAddress, `SELECT id, address FROM public.addresses LIMIT 1`)
			Expect(getErr).NotTo(HaveOccurred())
			expectedAddress := dbAddress{Id: addressId, Address: address}
			Expect(actualAddress).To(Equal(expectedAddress))
		})

		It("returns the existing record id if the address already exists", func() {
			_, createErr := repo.GetOrCreateAddressInTransaction(tx, address)
			Expect(createErr).NotTo(HaveOccurred())

			_, getErr := repo.GetOrCreateAddressInTransaction(tx, address)
			Expect(getErr).NotTo(HaveOccurred())
			tx.Commit()

			var addressCount int
			addressErr := db.Get(&addressCount, `SELECT count(*) FROM public.addresses`)
			Expect(addressErr).NotTo(HaveOccurred())
		})

		It("gets upper-cased addresses", func() {
			upperAddress := strings.ToUpper(address)
			upperAddressId, createErr := repo.GetOrCreateAddressInTransaction(tx, upperAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repo.GetOrCreateAddressInTransaction(tx, address)
			Expect(getErr).NotTo(HaveOccurred())
			tx.Commit()

			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})

		It("gets lower-cased addresses", func() {
			lowerAddress := strings.ToLower(address)
			upperAddressId, createErr := repo.GetOrCreateAddressInTransaction(tx, lowerAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repo.GetOrCreateAddressInTransaction(tx, address)
			Expect(getErr).NotTo(HaveOccurred())
			tx.Commit()

			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})
	})

	Describe("GetAddressById", func() {
		It("gets and address by it's id", func() {
			addressId, createErr := repo.GetOrCreateAddress(db, address)
			Expect(createErr).NotTo(HaveOccurred())

			actualAddress, getErr := repo.GetAddressById(db, addressId)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(actualAddress).To(Equal(address))
		})

		It("returns an error if the id doesn't exist", func() {
			_, getErr := repo.GetAddressById(db, 0)
			Expect(getErr).To(HaveOccurred())
			Expect(getErr).To(MatchError("sql: no rows in result set"))
		})
	})
})
