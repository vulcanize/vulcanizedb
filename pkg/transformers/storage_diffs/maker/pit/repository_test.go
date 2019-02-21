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
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/pit"
	. "github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
	"strconv"
)

var _ = Describe("Pit storage repository", func() {
	var (
		db              *postgres.DB
		err             error
		repo            pit.PitStorageRepository
		fakeAddress     = "0x12345"
		fakeBlockHash   = "expected_block_hash"
		fakeBlockNumber = 123
		fakeIlk         = "fake_ilk"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = pit.PitStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Ilk", func() {
		Describe("Line", func() {
			It("writes a row", func() {
				ilkLineMetadata := shared.GetStorageValueMetadata(pit.IlkLine, map[shared.Key]string{shared.Ilk: fakeIlk}, shared.Uint256)

				err = repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, line AS value FROM maker.pit_ilk_line`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkLineMetadata := shared.GetStorageValueMetadata(pit.IlkLine, nil, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLineMetadata, fakeUint256)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
			})
		})

		Describe("Spot", func() {
			It("writes a row", func() {
				ilkSpotMetadata := shared.GetStorageValueMetadata(pit.IlkSpot, map[shared.Key]string{shared.Ilk: fakeIlk}, shared.Uint256)

				err = repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, spot AS value FROM maker.pit_ilk_spot`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkSpotMetadata := shared.GetStorageValueMetadata(pit.IlkSpot, nil, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkSpotMetadata, fakeUint256)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
			})
		})
	})

	It("persists a pit drip", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, pit.DripMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, drip AS value FROM maker.pit_drip`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("persists a pit line", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, pit.LineMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, line AS value FROM maker.pit_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("persists a pit live", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, pit.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.pit_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("persists a pit vat", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, pit.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.pit_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})
})
