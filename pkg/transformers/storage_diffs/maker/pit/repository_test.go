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

package pit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/pit"
	. "github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pit storage repository", func() {
	var (
		blockNumber int
		blockHash   string
		db          *postgres.DB
		err         error
		repo        pit.PitStorageRepository
	)

	BeforeEach(func() {
		blockNumber = 123
		blockHash = "expected_block_hash"
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = pit.PitStorageRepository{}
		repo.SetDB(db)
	})

	It("persists an ilk line", func() {
		expectedIlk := "fake_ilk"
		expectedLine := "12345"
		ilkLineMetadata := shared.StorageValueMetadata{
			Name: pit.IlkLine,
			Keys: map[shared.Key]string{shared.Ilk: expectedIlk},
			Type: shared.Uint256,
		}
		err = repo.Create(blockNumber, blockHash, ilkLineMetadata, expectedLine)

		Expect(err).NotTo(HaveOccurred())

		var result MappingRes
		err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, line AS value FROM maker.pit_ilk_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertMapping(result, blockNumber, blockHash, expectedIlk, expectedLine)
	})

	It("persists an ilk spot", func() {
		expectedIlk := "fake_ilk"
		expectedSpot := "12345"
		ilkSpotMetadata := shared.StorageValueMetadata{
			Name: pit.IlkSpot,
			Keys: map[shared.Key]string{shared.Ilk: expectedIlk},
			Type: shared.Uint256,
		}
		err = repo.Create(blockNumber, blockHash, ilkSpotMetadata, expectedSpot)

		Expect(err).NotTo(HaveOccurred())

		var result MappingRes
		err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, spot AS value FROM maker.pit_ilk_spot`)
		Expect(err).NotTo(HaveOccurred())
		AssertMapping(result, blockNumber, blockHash, expectedIlk, expectedSpot)
	})

	It("persists a pit drip", func() {
		expectedDrip := "0x0123456789abcdef0123"

		err = repo.Create(blockNumber, blockHash, pit.DripMetadata, expectedDrip)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, drip AS value FROM maker.pit_drip`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedDrip)
	})

	It("persists a pit line", func() {
		expectedLine := "12345"

		err = repo.Create(blockNumber, blockHash, pit.LineMetadata, expectedLine)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, line AS value FROM maker.pit_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedLine)
	})

	It("persists a pit live", func() {
		expectedLive := "12345"

		err = repo.Create(blockNumber, blockHash, pit.LiveMetadata, expectedLive)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.pit_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedLive)
	})

	It("persists a pit vat", func() {
		expectedVat := "0x0123456789abcdef0123"

		err = repo.Create(blockNumber, blockHash, pit.VatMetadata, expectedVat)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.pit_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedVat)
	})
})
