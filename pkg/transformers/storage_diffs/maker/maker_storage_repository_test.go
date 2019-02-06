// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package maker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Maker storage repository", func() {
	It("fetches unique ilks from vat init events", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		insertVatInit("ilk1", 1, db)
		insertVatInit("ilk2", 2, db)
		insertVatInit("ilk2", 3, db)
		repository := maker.MakerStorageRepository{}
		repository.SetDB(db)

		ilks, err := repository.GetIlks()

		Expect(err).NotTo(HaveOccurred())
		Expect(len(ilks)).To(Equal(2))
		Expect(ilks).To(ConsistOf("ilk1", "ilk2"))
	})
})

func insertVatInit(ilk string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_init (header_id, ilk, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, ilk, 0, 0, "[]",
	)
	Expect(execErr).NotTo(HaveOccurred())
}
