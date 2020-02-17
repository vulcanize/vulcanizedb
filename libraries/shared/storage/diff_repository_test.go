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
	"database/sql"
	"math/rand"

	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage diffs repository", func() {
	var (
		db              = test_config.NewTestDB(test_config.NewTestNode())
		repo            storage.DiffRepository
		fakeStorageDiff types.RawDiff
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = storage.NewDiffRepository(db)
		fakeStorageDiff = types.RawDiff{
			HashedAddress: test_data.FakeHash(),
			BlockHash:     test_data.FakeHash(),
			BlockHeight:   rand.Int(),
			StorageKey:    test_data.FakeHash(),
			StorageValue:  test_data.FakeHash(),
		}
	})

	Describe("CreateStorageDiff", func() {
		It("adds a storage diff to the db, returning id", func() {
			id, createErr := repo.CreateStorageDiff(fakeStorageDiff)

			Expect(createErr).NotTo(HaveOccurred())
			Expect(id).NotTo(BeZero())
			var persisted types.PersistedDiff
			getErr := db.Get(&persisted, `SELECT * FROM public.storage_diff`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(persisted.ID).To(Equal(id))
			Expect(persisted.HashedAddress).To(Equal(fakeStorageDiff.HashedAddress))
			Expect(persisted.BlockHash).To(Equal(fakeStorageDiff.BlockHash))
			Expect(persisted.BlockHeight).To(Equal(fakeStorageDiff.BlockHeight))
			Expect(persisted.StorageKey).To(Equal(fakeStorageDiff.StorageKey))
			Expect(persisted.StorageValue).To(Equal(fakeStorageDiff.StorageValue))
		})

		It("does not duplicate storage diffs", func() {
			_, createErr := repo.CreateStorageDiff(fakeStorageDiff)
			Expect(createErr).NotTo(HaveOccurred())

			_, createTwoErr := repo.CreateStorageDiff(fakeStorageDiff)
			Expect(createTwoErr).To(HaveOccurred())
			Expect(createTwoErr).To(MatchError(sql.ErrNoRows))

			var count int
			getErr := db.Get(&count, `SELECT count(*) FROM public.storage_diff`)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("GetNewDiffs", func() {
		It("sends diffs that are not marked as checked", func() {
			diffs := make(chan types.PersistedDiff)
			errs := make(chan error)
			done := make(chan bool)
			fakeRawDiff := types.RawDiff{
				HashedAddress: test_data.FakeHash(),
				BlockHash:     test_data.FakeHash(),
				BlockHeight:   rand.Int(),
				StorageKey:    test_data.FakeHash(),
				StorageValue:  test_data.FakeHash(),
			}
			fakePersistedDiff := types.PersistedDiff{
				RawDiff: fakeRawDiff,
				ID:      rand.Int63(),
			}
			_, insertErr := db.Exec(`INSERT INTO public.storage_diff (id, block_height, block_hash,
				hashed_address, storage_key, storage_value) VALUES ($1, $2, $3, $4, $5, $6)`, fakePersistedDiff.ID,
				fakeRawDiff.BlockHeight, fakeRawDiff.BlockHash.Bytes(), fakeRawDiff.HashedAddress.Bytes(),
				fakeRawDiff.StorageKey.Bytes(), fakeRawDiff.StorageValue.Bytes())
			Expect(insertErr).NotTo(HaveOccurred())

			go repo.GetNewDiffs(diffs, errs, done)

			Consistently(errs).ShouldNot(Receive())
			Eventually(<-diffs).Should(Equal(fakePersistedDiff))
			Eventually(<-done).Should(BeTrue())
		})

		It("does not send diff that's marked as checked", func() {
			diffs := make(chan types.PersistedDiff)
			errs := make(chan error)
			done := make(chan bool)
			fakeRawDiff := types.RawDiff{
				HashedAddress: test_data.FakeHash(),
				BlockHash:     test_data.FakeHash(),
				BlockHeight:   rand.Int(),
				StorageKey:    test_data.FakeHash(),
				StorageValue:  test_data.FakeHash(),
			}
			fakePersistedDiff := types.PersistedDiff{
				RawDiff: fakeRawDiff,
				ID:      rand.Int63(),
				Checked: true,
			}
			_, insertErr := db.Exec(`INSERT INTO public.storage_diff (id, block_height, block_hash,
				hashed_address, storage_key, storage_value, checked) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				fakePersistedDiff.ID, fakeRawDiff.BlockHeight, fakeRawDiff.BlockHash.Bytes(),
				fakeRawDiff.HashedAddress.Bytes(), fakeRawDiff.StorageKey.Bytes(), fakeRawDiff.StorageValue.Bytes(),
				fakePersistedDiff.Checked)
			Expect(insertErr).NotTo(HaveOccurred())

			go repo.GetNewDiffs(diffs, errs, done)

			Consistently(errs).ShouldNot(Receive())
			Consistently(diffs).ShouldNot(Receive())
			Eventually(<-done).Should(BeTrue())
		})
	})

	Describe("MarkChecked", func() {
		It("marks a diff as checked", func() {
			fakeRawDiff := types.RawDiff{
				HashedAddress: test_data.FakeHash(),
				BlockHash:     test_data.FakeHash(),
				BlockHeight:   rand.Int(),
				StorageKey:    test_data.FakeHash(),
				StorageValue:  test_data.FakeHash(),
			}
			fakePersistedDiff := types.PersistedDiff{
				RawDiff: fakeRawDiff,
				ID:      rand.Int63(),
				Checked: true,
			}
			_, insertErr := db.Exec(`INSERT INTO public.storage_diff (id, block_height, block_hash,
				hashed_address, storage_key, storage_value, checked) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				fakePersistedDiff.ID, fakeRawDiff.BlockHeight, fakeRawDiff.BlockHash.Bytes(),
				fakeRawDiff.HashedAddress.Bytes(), fakeRawDiff.StorageKey.Bytes(), fakeRawDiff.StorageValue.Bytes(),
				fakePersistedDiff.Checked)
			Expect(insertErr).NotTo(HaveOccurred())

			err := repo.MarkChecked(fakePersistedDiff.ID)

			Expect(err).NotTo(HaveOccurred())
			var checked bool
			checkedErr := db.Get(&checked, `SELECT checked FROM public.storage_diff WHERE id = $1`, fakePersistedDiff.ID)
			Expect(checkedErr).NotTo(HaveOccurred())
			Expect(checked).To(BeTrue())
		})
	})

	Describe("MarkFromBackfill", func() {
		It("marks a diff as from_backfill", func() {
			id, createErr := repo.CreateStorageDiff(fakeStorageDiff)
			Expect(createErr).NotTo(HaveOccurred())

			err := repo.MarkFromBackfill(id)

			Expect(err).NotTo(HaveOccurred())
			var fromBackfill bool
			checkedErr := db.Get(&fromBackfill, `SELECT from_backfill FROM public.storage_diff WHERE id = $1`, id)
			Expect(checkedErr).NotTo(HaveOccurred())
			Expect(fromBackfill).To(BeTrue())
		})
	})
})
