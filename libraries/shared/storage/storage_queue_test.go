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
		diff  utils.StorageDiff
		queue storage.IStorageQueue
	)

	BeforeEach(func() {
		fakeAddr := "0x123456"
		diff = utils.StorageDiff{
			HashedAddress: utils.HexToKeccak256Hash(fakeAddr),
			BlockHash:     common.HexToHash("0x678901"),
			BlockHeight:   987,
			StorageKey:    common.HexToHash("0x654321"),
			StorageValue:  common.HexToHash("0x198765"),
		}
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		queue = storage.NewStorageQueue(db)
		addErr := queue.Add(diff)
		Expect(addErr).NotTo(HaveOccurred())
	})

	Describe("Add", func() {
		It("adds a storage diff to the db", func() {
			var result utils.StorageDiff
			getErr := db.Get(&result, `SELECT contract, block_hash, block_height, storage_key, storage_value FROM public.queued_storage`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(diff))
		})

		It("does not duplicate storage diffs", func() {
			addErr := queue.Add(diff)
			Expect(addErr).NotTo(HaveOccurred())
			var count int
			getErr := db.Get(&count, `SELECT count(*) FROM public.queued_storage`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	It("deletes storage diff from db", func() {
		diffs, getErr := queue.GetAll()
		Expect(getErr).NotTo(HaveOccurred())
		Expect(len(diffs)).To(Equal(1))

		err := queue.Delete(diffs[0].Id)

		Expect(err).NotTo(HaveOccurred())
		remainingRows, secondGetErr := queue.GetAll()
		Expect(secondGetErr).NotTo(HaveOccurred())
		Expect(len(remainingRows)).To(BeZero())
	})

	It("gets all storage diffs from db", func() {
		fakeAddr := "0x234567"
		diffTwo := utils.StorageDiff{
			HashedAddress: utils.HexToKeccak256Hash(fakeAddr),
			BlockHash:     common.HexToHash("0x678902"),
			BlockHeight:   988,
			StorageKey:    common.HexToHash("0x654322"),
			StorageValue:  common.HexToHash("0x198766"),
		}
		addErr := queue.Add(diffTwo)
		Expect(addErr).NotTo(HaveOccurred())

		diffs, err := queue.GetAll()

		Expect(err).NotTo(HaveOccurred())
		Expect(len(diffs)).To(Equal(2))
		Expect(diffs[0]).NotTo(Equal(diffs[1]))
		Expect(diffs[0].Id).NotTo(BeZero())
		Expect(diffs[0].HashedAddress).To(Or(Equal(diff.HashedAddress), Equal(diffTwo.HashedAddress)))
		Expect(diffs[0].BlockHash).To(Or(Equal(diff.BlockHash), Equal(diffTwo.BlockHash)))
		Expect(diffs[0].BlockHeight).To(Or(Equal(diff.BlockHeight), Equal(diffTwo.BlockHeight)))
		Expect(diffs[0].StorageKey).To(Or(Equal(diff.StorageKey), Equal(diffTwo.StorageKey)))
		Expect(diffs[0].StorageValue).To(Or(Equal(diff.StorageValue), Equal(diffTwo.StorageValue)))
		Expect(diffs[1].Id).NotTo(BeZero())
		Expect(diffs[1].HashedAddress).To(Or(Equal(diff.HashedAddress), Equal(diffTwo.HashedAddress)))
		Expect(diffs[1].BlockHash).To(Or(Equal(diff.BlockHash), Equal(diffTwo.BlockHash)))
		Expect(diffs[1].BlockHeight).To(Or(Equal(diff.BlockHeight), Equal(diffTwo.BlockHeight)))
		Expect(diffs[1].StorageKey).To(Or(Equal(diff.StorageKey), Equal(diffTwo.StorageKey)))
		Expect(diffs[1].StorageValue).To(Or(Equal(diff.StorageValue), Equal(diffTwo.StorageValue)))
	})
})
