/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/vow"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vow storage repository test", func() {
	var (
		blockNumber int
		blockHash   string
		db          *postgres.DB
		err         error
		repo        vow.VowStorageRepository
	)

	BeforeEach(func() {
		blockNumber = 123
		blockHash = "expected_block_hash"
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vow.VowStorageRepository{}
		repo.SetDB(db)
	})

	It("persists a vow vat", func() {
		expectedVat := "123"

		err = repo.Create(blockNumber, blockHash, vow.VatMetadata, expectedVat)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value from maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedVat)
	})

	It("persists a vow cow", func() {
		expectedCow := "123"

		err = repo.Create(blockNumber, blockHash, vow.CowMetadata, expectedCow)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, cow AS value from maker.vow_cow`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedCow)
	})

	It("persists a vow row", func() {
		expectedRow := "123"

		err = repo.Create(blockNumber, blockHash, vow.RowMetadata, expectedRow)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, row AS value from maker.vow_row`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedRow)
	})

	It("persists a vow Sin", func() {
		expectedSow := "123"

		err = repo.Create(blockNumber, blockHash, vow.SinMetadata, expectedSow)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, sin AS value from maker.vow_sin`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedSow)
	})

	It("persists a vow woe", func() {
		expectedWoe := "123"

		err = repo.Create(blockNumber, blockHash, vow.WoeMetadata, expectedWoe)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, woe AS value from maker.vow_woe`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedWoe)
	})

	It("persists a vow Ash", func() {
		expectedAsh := "123"

		err = repo.Create(blockNumber, blockHash, vow.AshMetadata, expectedAsh)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, ash AS value from maker.vow_ash`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedAsh)
	})

	It("persists a vow Wait", func() {
		expectedWait := "123"

		err = repo.Create(blockNumber, blockHash, vow.WaitMetadata, expectedWait)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, wait AS value from maker.vow_wait`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedWait)
	})

	It("persists a vow Bump", func() {
		expectedBump := "123"

		err = repo.Create(blockNumber, blockHash, vow.BumpMetadata, expectedBump)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, bump AS value from maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedBump)
	})

	It("persists a vow Sump", func() {
		expectedSump := "123"

		err = repo.Create(blockNumber, blockHash, vow.SumpMetadata, expectedSump)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, sump AS value from maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedSump)
	})

	It("persists a vow Hump", func() {
		expectedHump := "123"

		err = repo.Create(blockNumber, blockHash, vow.HumpMetadata, expectedHump)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, hump AS value from maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedHump)
	})
})
