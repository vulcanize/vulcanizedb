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
	. "github.com/onsi/ginkgo"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
	"strings"

	. "github.com/onsi/gomega"
)

var _ = Describe("address repository", func() {
	var (
		db *postgres.DB
		repo repositories.AddressRepository
		address = fakes.FakeAddress.Hex()
	)
	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = repositories.AddressRepository{}
	})

	type dbAddress struct {
		Id int
		Address string
	}

	It("creates an address record", func() {
		addressId, createErr := repo.CreateOrGetAddress(db, address)
		Expect(createErr).NotTo(HaveOccurred())

		var actualAddress dbAddress
		getErr := db.Get(&actualAddress, `SELECT id, address FROM public.addresses LIMIT 1`)
		Expect(getErr).NotTo(HaveOccurred())
		expectedAddress := dbAddress{Id: addressId, Address: address}
		Expect(actualAddress).To(Equal(expectedAddress))
	})

	It("returns the existing record id if the address already exists", func() {
		_, createErr := repo.CreateOrGetAddress(db, address)
		Expect(createErr).NotTo(HaveOccurred())

		_, getErr := repo.CreateOrGetAddress(db, address)
		Expect(getErr).NotTo(HaveOccurred())

		var addressCount int
		addressErr := db.Get(&addressCount, `SELECT count(*) FROM public.addresses`)
		Expect(addressErr).NotTo(HaveOccurred())
	})

	It("gets upper-cased addresses", func() {
		//insert it as all upper
		upperAddress := strings.ToUpper(address)
		upperAddressId, createErr := repo.CreateOrGetAddress(db, upperAddress)
		Expect(createErr).NotTo(HaveOccurred())

		mixedCaseAddressId, getErr := repo.CreateOrGetAddress(db, address)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(upperAddressId).To(Equal(mixedCaseAddressId))
	})

	It("gets lower-cased addresses", func() {
		//insert it as all upper
		lowerAddress := strings.ToLower(address)
		upperAddressId, createErr := repo.CreateOrGetAddress(db, lowerAddress)
		Expect(createErr).NotTo(HaveOccurred())

		mixedCaseAddressId, getErr := repo.CreateOrGetAddress(db, address)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(upperAddressId).To(Equal(mixedCaseAddressId))
	})
})