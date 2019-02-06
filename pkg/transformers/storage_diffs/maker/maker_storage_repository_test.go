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
	var (
		db         *postgres.DB
		repository maker.IMakerStorageRepository
		ilk1       = "ilk1"
		ilk2       = "ilk2"
		guy1       = "guy1"
		guy2       = "guy2"
		guy3       = "guy3"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = &maker.MakerStorageRepository{}
		repository.SetDB(db)
	})

	Describe("getting dai keys", func() {
		It("fetches guy from both src and dst field on vat_move", func() {
			insertVatMove(guy1, guy2, 1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(2))
			Expect(keys).To(ConsistOf(guy1, guy2))
		})

		It("fetches guy from w field on vat_tune", func() {
			insertVatTune(guy1, guy1, guy1, guy2, 1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy2))
		})

		It("fetches guy from v field on vat_heal", func() {
			insertVatHeal(guy2, guy1, 1, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(1))
			Expect(keys).To(ConsistOf(guy1))
		})

		It("fetches unique guys from vat_move + vat_tune + vat_heal + vat_fold", func() {
			guy4 := "guy4"
			guy5 := "guy5"
			guy6 := "guy6"
			insertVatMove(guy1, guy2, 1, db)
			insertVatTune(guy1, guy1, guy1, guy3, 2, db)
			insertVatHeal(guy6, guy4, 3, db)
			insertVatFold(guy5, 4, db)
			// duplicates
			insertVatMove(guy3, guy1, 5, db)
			insertVatTune(guy2, guy2, guy2, guy5, 6, db)
			insertVatHeal(guy6, guy2, 7, db)
			insertVatFold(guy4, 8, db)

			keys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(keys)).To(Equal(5))
			Expect(keys).To(ConsistOf(guy1, guy2, guy3, guy4, guy5))
		})

		It("does not return error if no matching rows", func() {
			daiKeys, err := repository.GetDaiKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(daiKeys)).To(BeZero())
		})
	})

	Describe("getting gem keys", func() {
		It("fetches guy from both src and dst field on vat_flux", func() {
			insertVatFlux(ilk1, guy1, guy2, 1, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(2))
			Expect(gems).To(ConsistOf([]maker.Urn{{
				Ilk: ilk1,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy2,
			}}))
		})

		It("fetches guy from v field on vat_tune + vat_grab", func() {
			insertVatTune(ilk1, guy1, guy2, guy1, 1, db)
			insertVatGrab(ilk1, guy1, guy3, guy1, 2, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(2))
			Expect(gems).To(ConsistOf([]maker.Urn{{
				Ilk: ilk1,
				Guy: guy2,
			}, {
				Ilk: ilk1,
				Guy: guy3,
			}}))
		})

		It("fetches unique urns from vat_slip + vat_flux + vat_tune + vat_grab + vat_toll events", func() {
			insertVatSlip(ilk1, guy1, 1, db)
			insertVatFlux(ilk1, guy2, guy3, 2, db)
			insertVatTune(ilk2, guy1, guy1, guy1, 3, db)
			insertVatGrab(ilk2, guy1, guy2, guy1, 4, db)
			insertVatToll(ilk2, guy3, 5, db)
			// duplicates
			insertVatSlip(ilk1, guy2, 6, db)
			insertVatFlux(ilk2, guy2, guy3, 7, db)
			insertVatTune(ilk2, guy1, guy1, guy1, 8, db)
			insertVatGrab(ilk1, guy1, guy1, guy1, 9, db)
			insertVatToll(ilk1, guy3, 10, db)

			gems, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gems)).To(Equal(6))
			Expect(gems).To(ConsistOf([]maker.Urn{{
				Ilk: ilk1,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy2,
			}, {
				Ilk: ilk1,
				Guy: guy3,
			}, {
				Ilk: ilk2,
				Guy: guy1,
			}, {
				Ilk: ilk2,
				Guy: guy2,
			}, {
				Ilk: ilk2,
				Guy: guy3,
			}}))
		})

		It("does not return error if no matching rows", func() {
			gemKeys, err := repository.GetGemKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(gemKeys)).To(BeZero())
		})
	})

	Describe("getting ilks", func() {
		It("fetches unique ilks from vat init events", func() {
			insertVatInit(ilk1, 1, db)
			insertVatInit(ilk2, 2, db)
			insertVatInit(ilk2, 3, db)

			ilks, err := repository.GetIlks()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(ilks)).To(Equal(2))
			Expect(ilks).To(ConsistOf(ilk1, ilk2))
		})

		It("does not return error if no matching rows", func() {
			ilks, err := repository.GetIlks()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(ilks)).To(BeZero())
		})
	})

	Describe("getting sin keys", func() {
		It("fetches guy from w field of vat grab", func() {
			insertVatGrab(guy1, guy1, guy1, guy2, 1, db)

			sinKeys, err := repository.GetSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy2))
		})

		It("fetches guy from u field of vat heal", func() {
			insertVatHeal(guy1, guy2, 1, db)

			sinKeys, err := repository.GetSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(1))
			Expect(sinKeys).To(ConsistOf(guy1))
		})

		It("fetches unique sin keys from vat_grab + vat_heal", func() {
			insertVatGrab(guy3, guy3, guy3, guy1, 1, db)
			insertVatHeal(guy2, guy3, 2, db)
			// duplicates
			insertVatGrab(guy2, guy2, guy2, guy2, 3, db)
			insertVatHeal(guy1, guy2, 4, db)

			sinKeys, err := repository.GetSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(Equal(2))
			Expect(sinKeys).To(ConsistOf(guy1, guy2))
		})

		It("does not return error if no matching rows", func() {
			sinKeys, err := repository.GetSinKeys()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(sinKeys)).To(BeZero())
		})
	})

	Describe("getting urns", func() {
		It("fetches unique urns from vat_tune + vat_grab events", func() {
			insertVatTune(ilk1, guy1, guy1, guy1, 1, db)
			insertVatTune(ilk1, guy2, guy1, guy1, 2, db)
			insertVatTune(ilk2, guy1, guy1, guy1, 3, db)
			insertVatTune(ilk1, guy1, guy1, guy1, 4, db)
			insertVatGrab(ilk1, guy1, guy1, guy1, 5, db)
			insertVatGrab(ilk1, guy3, guy1, guy1, 6, db)

			urns, err := repository.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(Equal(4))
			Expect(urns).To(ConsistOf([]maker.Urn{{
				Ilk: ilk1,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy2,
			}, {
				Ilk: ilk2,
				Guy: guy1,
			}, {
				Ilk: ilk1,
				Guy: guy3,
			}}))
		})

		It("does not return error if no matching rows", func() {
			urns, err := repository.GetUrns()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(urns)).To(BeZero())
		})
	})
})

func insertVatFold(urn string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_fold (header_id, urn, log_idx, tx_idx)
			VALUES($1, $2, $3, $4)`,
		headerID, urn, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatFlux(ilk, src, dst string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_flux (header_id, ilk, src, dst, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, ilk, src, dst, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatGrab(ilk, urn, v, w string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_grab (header_id, ilk, urn, v, w, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
		headerID, ilk, urn, v, w, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatHeal(urn, v string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_heal (header_id, urn, v, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, urn, v, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatInit(ilk string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_init (header_id, ilk, log_idx, tx_idx)
			VALUES($1, $2, $3, $4)`,
		headerID, ilk, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatMove(src, dst string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_move (header_id, src, dst, rad, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6)`,
		headerID, src, dst, 0, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatSlip(ilk, guy string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_slip (header_id, ilk, guy, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, ilk, guy, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatToll(ilk, urn string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_toll (header_id, ilk, urn, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5)`,
		headerID, ilk, urn, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}

func insertVatTune(ilk, urn, v, w string, blockNumber int64, db *postgres.DB) {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
	Expect(err).NotTo(HaveOccurred())
	_, execErr := db.Exec(
		`INSERT INTO maker.vat_tune (header_id, ilk, urn, v, w, log_idx, tx_idx)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
		headerID, ilk, urn, v, w, 0, 0,
	)
	Expect(execErr).NotTo(HaveOccurred())
}
