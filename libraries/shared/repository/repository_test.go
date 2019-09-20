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

package repository_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("", func() {
	Describe("MarkContractWatcherHeaderCheckedInTransaction", func() {
		var (
			checkedHeadersColumn string
			db                   *postgres.DB
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			checkedHeadersColumn = "test_column_checked"
			_, migrateErr := db.Exec(`ALTER TABLE public.checked_headers
				ADD COLUMN ` + checkedHeadersColumn + ` integer`)
			Expect(migrateErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, cleanupMigrateErr := db.Exec(`ALTER TABLE public.checked_headers DROP COLUMN ` + checkedHeadersColumn)
			Expect(cleanupMigrateErr).NotTo(HaveOccurred())
		})

		It("marks passed header as checked within a passed transaction", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			err := repository.MarkContractWatcherHeaderCheckedInTransaction(headerID, tx, checkedHeadersColumn)
			Expect(err).NotTo(HaveOccurred())
			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())
			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT COUNT(*) FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(1))
		})
	})
})
