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
		type IlkLine struct {
			BlockMetadata
			Ilk  string
			Line string
		}
		var result IlkLine
		err = db.Get(&result, `SELECT block_number, block_hash, ilk, line FROM maker.pit_ilk_line`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Ilk).To(Equal(expectedIlk))
		Expect(result.Line).To(Equal(expectedLine))
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
		type IlkSpot struct {
			BlockMetadata
			Ilk  string
			Spot string
		}
		var result IlkSpot
		err = db.Get(&result, `SELECT block_number, block_hash, ilk, spot FROM maker.pit_ilk_spot`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Ilk).To(Equal(expectedIlk))
		Expect(result.Spot).To(Equal(expectedSpot))
	})

	It("persists a pit drip", func() {
		expectedDrip := "0x0123456789abcdef0123"

		err = repo.Create(blockNumber, blockHash, pit.DripMetadata, expectedDrip)

		Expect(err).NotTo(HaveOccurred())
		type PitDrip struct {
			BlockMetadata
			Drip string
		}
		var result PitDrip
		err = db.Get(&result, `SELECT block_number, block_hash, drip FROM maker.pit_drip`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Drip).To(Equal(expectedDrip))
	})

	It("persists a pit line", func() {
		expectedLine := "12345"

		err = repo.Create(blockNumber, blockHash, pit.LineMetadata, expectedLine)

		Expect(err).NotTo(HaveOccurred())
		type PitLine struct {
			BlockMetadata
			Line string
		}
		var result PitLine
		err = db.Get(&result, `SELECT block_number, block_hash, line FROM maker.pit_line`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Line).To(Equal(expectedLine))
	})

	It("persists a pit live", func() {
		expectedLive := "12345"

		err = repo.Create(blockNumber, blockHash, pit.LiveMetadata, expectedLive)

		Expect(err).NotTo(HaveOccurred())
		type PitLive struct {
			BlockMetadata
			Live string
		}
		var result PitLive
		err = db.Get(&result, `SELECT block_number, block_hash, live FROM maker.pit_live`)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.BlockNumber).To(Equal(blockNumber))
		Expect(result.BlockHash).To(Equal(blockHash))
		Expect(result.Live).To(Equal(expectedLive))
	})

	It("persists a pit vat", func() {
		expectedVat := "0x0123456789abcdef0123"

		err = repo.Create(blockNumber, blockHash, pit.VatMetadata, expectedVat)

		Expect(err).NotTo(HaveOccurred())
		type PitVat struct {
			BlockMetadata
			Vat string
		}
		var result PitVat
		err = db.Get(&result, `SELECT block_number, block_hash, vat FROM maker.pit_vat`)
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
