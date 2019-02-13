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
		type VowVat struct {
			BlockMetadata
			Vat string
		}
		var result VowVat
		err = db.Get(&result, `SELECT block_number, block_hash, vat from maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Vat).To(Equal(expectedVat))
	})

	It("persists a vow cow", func() {
		expectedCow := "123"

		err = repo.Create(blockNumber, blockHash, vow.CowMetadata, expectedCow)

		Expect(err).NotTo(HaveOccurred())
		type VowCow struct {
			BlockMetadata
			Cow string
		}
		var result VowCow
		err = db.Get(&result, `SELECT block_number, block_hash, cow from maker.vow_cow`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Cow).To(Equal(expectedCow))
	})

	It("persists a vow row", func() {
		expectedRow := "123"

		err = repo.Create(blockNumber, blockHash, vow.RowMetadata, expectedRow)

		Expect(err).NotTo(HaveOccurred())
		type VowRow struct {
			BlockMetadata
			Row string
		}
		var result VowRow
		err = db.Get(&result, `SELECT block_number, block_hash, row from maker.vow_row`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Row).To(Equal(expectedRow))
	})

	It("persists a vow Sin", func() {
		expectedSow := "123"

		err = repo.Create(blockNumber, blockHash, vow.SinMetadata, expectedSow)

		Expect(err).NotTo(HaveOccurred())
		type VowSin struct {
			BlockMetadata
			Sin string
		}
		var result VowSin
		err = db.Get(&result, `SELECT block_number, block_hash, sin from maker.vow_sin`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Sin).To(Equal(expectedSow))
	})

	It("persists a vow woe", func() {
		expectedWoe := "123"

		err = repo.Create(blockNumber, blockHash, vow.WoeMetadata, expectedWoe)

		Expect(err).NotTo(HaveOccurred())
		type VowWoe struct {
			BlockMetadata
			Woe string
		}
		var result VowWoe
		err = db.Get(&result, `SELECT block_number, block_hash, woe from maker.vow_woe`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Woe).To(Equal(expectedWoe))
	})

	It("persists a vow Ash", func() {
		expectedAsh := "123"

		err = repo.Create(blockNumber, blockHash, vow.AshMetadata, expectedAsh)

		Expect(err).NotTo(HaveOccurred())
		type VowAsh struct {
			BlockMetadata
			Ash string
		}
		var result VowAsh
		err = db.Get(&result, `SELECT block_number, block_hash, ash from maker.vow_ash`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Ash).To(Equal(expectedAsh))
	})

	It("persists a vow Wait", func() {
		expectedWait := "123"

		err = repo.Create(blockNumber, blockHash, vow.WaitMetadata, expectedWait)

		Expect(err).NotTo(HaveOccurred())
		type VowWait struct {
			BlockMetadata
			Wait string
		}
		var result VowWait
		err = db.Get(&result, `SELECT block_number, block_hash, wait from maker.vow_wait`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Wait).To(Equal(expectedWait))
	})

	It("persists a vow Bump", func() {
		expectedBump := "123"

		err = repo.Create(blockNumber, blockHash, vow.BumpMetadata, expectedBump)

		Expect(err).NotTo(HaveOccurred())
		type VowBump struct {
			BlockMetadata
			Bump string
		}
		var result VowBump
		err = db.Get(&result, `SELECT block_number, block_hash, bump from maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Bump).To(Equal(expectedBump))
	})

	It("persists a vow Sump", func() {
		expectedSump := "123"

		err = repo.Create(blockNumber, blockHash, vow.SumpMetadata, expectedSump)

		Expect(err).NotTo(HaveOccurred())
		type VowSump struct {
			BlockMetadata
			Sump string
		}
		var result VowSump
		err = db.Get(&result, `SELECT block_number, block_hash, sump from maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Sump).To(Equal(expectedSump))
	})

	It("persists a vow Hump", func() {
		expectedHump := "123"

		err = repo.Create(blockNumber, blockHash, vow.HumpMetadata, expectedHump)

		Expect(err).NotTo(HaveOccurred())
		type VowHump struct {
			BlockMetadata
			Hump string
		}
		var result VowHump
		err = db.Get(&result, `SELECT block_number, block_hash, hump from maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Hump).To(Equal(expectedHump))
	})
})

type BlockMetadata struct {
	BlockNumber int    `db:"block_number"`
	BlockHash   string `db:"block_hash"`
}
