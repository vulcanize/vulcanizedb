// VulcanizeDB
// Copyright © 2019 Vulcanize
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package repository_test

import (
	"strings"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("address lookup", func() {
	var (
		db      *postgres.DB
		address = fakes.FakeAddress.Hex()
	)
	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
	})

	AfterEach(func() {
		test_config.CleanTestDB(db)
	})

	type dbAddress struct {
		Id            int64
		Address       string
		HashedAddress string `db:"hashed_address"`
	}

	Describe("GetOrCreateAddress", func() {
		It("creates an address record", func() {
			addressId, createErr := repository.GetOrCreateAddress(db, address)
			Expect(createErr).NotTo(HaveOccurred())

			var actualAddress dbAddress
			getErr := db.Get(&actualAddress, `SELECT id, address, hashed_address FROM public.addresses LIMIT 1`)
			Expect(getErr).NotTo(HaveOccurred())
			hashedAddress := utils.HexToKeccak256Hash(address).Hex()
			expectedAddress := dbAddress{Id: addressId, Address: address, HashedAddress: hashedAddress}
			Expect(actualAddress).To(Equal(expectedAddress))
		})

		It("returns the existing record id if the address already exists", func() {
			createId, createErr := repository.GetOrCreateAddress(db, address)
			Expect(createErr).NotTo(HaveOccurred())

			getId, getErr := repository.GetOrCreateAddress(db, address)
			Expect(getErr).NotTo(HaveOccurred())

			var addressCount int
			addressErr := db.Get(&addressCount, `SELECT count(*) FROM public.addresses`)
			Expect(addressErr).NotTo(HaveOccurred())
			Expect(addressCount).To(Equal(1))
			Expect(createId).To(Equal(getId))
		})

		It("gets upper-cased addresses", func() {
			upperAddress := strings.ToUpper(address)
			upperAddressId, createErr := repository.GetOrCreateAddress(db, upperAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repository.GetOrCreateAddress(db, address)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})

		It("gets lower-cased addresses", func() {
			lowerAddress := strings.ToLower(address)
			upperAddressId, createErr := repository.GetOrCreateAddress(db, lowerAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repository.GetOrCreateAddress(db, address)
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

		AfterEach(func() {
			tx.Rollback()
		})

		It("creates an address record", func() {
			addressId, createErr := repository.GetOrCreateAddressInTransaction(tx, address)
			Expect(createErr).NotTo(HaveOccurred())
			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			var actualAddress dbAddress
			getErr := db.Get(&actualAddress, `SELECT id, address, hashed_address FROM public.addresses LIMIT 1`)
			Expect(getErr).NotTo(HaveOccurred())
			hashedAddress := utils.HexToKeccak256Hash(address).Hex()
			expectedAddress := dbAddress{Id: addressId, Address: address, HashedAddress: hashedAddress}
			Expect(actualAddress).To(Equal(expectedAddress))
		})

		It("returns the existing record id if the address already exists", func() {
			_, createErr := repository.GetOrCreateAddressInTransaction(tx, address)
			Expect(createErr).NotTo(HaveOccurred())

			_, getErr := repository.GetOrCreateAddressInTransaction(tx, address)
			Expect(getErr).NotTo(HaveOccurred())
			tx.Commit()

			var addressCount int
			addressErr := db.Get(&addressCount, `SELECT count(*) FROM public.addresses`)
			Expect(addressErr).NotTo(HaveOccurred())
		})

		It("gets upper-cased addresses", func() {
			upperAddress := strings.ToUpper(address)
			upperAddressId, createErr := repository.GetOrCreateAddressInTransaction(tx, upperAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repository.GetOrCreateAddressInTransaction(tx, address)
			Expect(getErr).NotTo(HaveOccurred())
			tx.Commit()

			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})

		It("gets lower-cased addresses", func() {
			lowerAddress := strings.ToLower(address)
			upperAddressId, createErr := repository.GetOrCreateAddressInTransaction(tx, lowerAddress)
			Expect(createErr).NotTo(HaveOccurred())

			mixedCaseAddressId, getErr := repository.GetOrCreateAddressInTransaction(tx, address)
			Expect(getErr).NotTo(HaveOccurred())
			tx.Commit()

			Expect(upperAddressId).To(Equal(mixedCaseAddressId))
		})
	})

	Describe("GetAddressByID", func() {
		It("gets and address by it's id", func() {
			addressId, createErr := repository.GetOrCreateAddress(db, address)
			Expect(createErr).NotTo(HaveOccurred())

			actualAddress, getErr := repository.GetAddressByID(db, addressId)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(actualAddress).To(Equal(address))
		})

		It("returns an error if the id doesn't exist", func() {
			_, getErr := repository.GetAddressByID(db, 0)
			Expect(getErr).To(HaveOccurred())
			Expect(getErr).To(MatchError("sql: no rows in result set"))
		})
	})
})
