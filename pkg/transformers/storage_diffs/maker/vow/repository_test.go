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
		expectedVat := "yo"

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
})

type BlockMetadata struct {
	BlockNumber int    `db:"block_number"`
	BlockHash   string `db:"block_hash"`
}
