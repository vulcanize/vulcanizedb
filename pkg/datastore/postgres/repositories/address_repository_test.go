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

package repositories_test

import (
	"strings"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
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
			Expect(addressCount).To(Equal(1))
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
