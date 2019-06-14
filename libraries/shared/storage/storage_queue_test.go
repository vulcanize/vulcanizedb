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

package storage_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage queue", func() {
	var (
		db    *postgres.DB
		row   utils.StorageDiffRow
		queue storage.IStorageQueue
	)

	BeforeEach(func() {
		row = utils.StorageDiffRow{
			Contract:     common.HexToAddress("0x123456"),
			BlockHash:    common.HexToHash("0x678901"),
			BlockHeight:  987,
			StorageKey:   common.HexToHash("0x654321"),
			StorageValue: common.HexToHash("0x198765"),
		}
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		queue = storage.NewStorageQueue(db)
		addErr := queue.Add(row)
		Expect(addErr).NotTo(HaveOccurred())
	})

	Describe("Add", func() {
		It("adds a storage row to the db", func() {
			var result utils.StorageDiffRow
			getErr := db.Get(&result, `SELECT contract, block_hash, block_height, storage_key, storage_value FROM public.queued_storage`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(row))
		})

		It("does not duplicate storage rows", func() {
			addErr := queue.Add(row)
			Expect(addErr).NotTo(HaveOccurred())
			var count int
			getErr := db.Get(&count, `SELECT count(*) FROM public.queued_storage`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	It("deletes storage row from db", func() {
		rows, getErr := queue.GetAll()
		Expect(getErr).NotTo(HaveOccurred())
		Expect(len(rows)).To(Equal(1))

		err := queue.Delete(rows[0].Id)

		Expect(err).NotTo(HaveOccurred())
		remainingRows, secondGetErr := queue.GetAll()
		Expect(secondGetErr).NotTo(HaveOccurred())
		Expect(len(remainingRows)).To(BeZero())
	})

	It("gets all storage rows from db", func() {
		rowTwo := utils.StorageDiffRow{
			Contract:     common.HexToAddress("0x123456"),
			BlockHash:    common.HexToHash("0x678902"),
			BlockHeight:  988,
			StorageKey:   common.HexToHash("0x654322"),
			StorageValue: common.HexToHash("0x198766"),
		}
		addErr := queue.Add(rowTwo)
		Expect(addErr).NotTo(HaveOccurred())

		rows, err := queue.GetAll()

		Expect(err).NotTo(HaveOccurred())
		Expect(len(rows)).To(Equal(2))
		Expect(rows[0]).NotTo(Equal(rows[1]))
		Expect(rows[0].Id).NotTo(BeZero())
		Expect(rows[0].Contract).To(Or(Equal(row.Contract), Equal(rowTwo.Contract)))
		Expect(rows[0].BlockHash).To(Or(Equal(row.BlockHash), Equal(rowTwo.BlockHash)))
		Expect(rows[0].BlockHeight).To(Or(Equal(row.BlockHeight), Equal(rowTwo.BlockHeight)))
		Expect(rows[0].StorageKey).To(Or(Equal(row.StorageKey), Equal(rowTwo.StorageKey)))
		Expect(rows[0].StorageValue).To(Or(Equal(row.StorageValue), Equal(rowTwo.StorageValue)))
		Expect(rows[1].Id).NotTo(BeZero())
		Expect(rows[1].Contract).To(Or(Equal(row.Contract), Equal(rowTwo.Contract)))
		Expect(rows[1].BlockHash).To(Or(Equal(row.BlockHash), Equal(rowTwo.BlockHash)))
		Expect(rows[1].BlockHeight).To(Or(Equal(row.BlockHeight), Equal(rowTwo.BlockHeight)))
		Expect(rows[1].StorageKey).To(Or(Equal(row.StorageKey), Equal(rowTwo.StorageKey)))
		Expect(rows[1].StorageValue).To(Or(Equal(row.StorageValue), Equal(rowTwo.StorageValue)))
	})
})
