package vat_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	. "github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/vat"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
	"strconv"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db              *postgres.DB
		repo            vat.VatStorageRepository
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeIlk         = "fake_ilk"
		fakeGuy         = "fake_urn"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vat.VatStorageRepository{}
		repo.SetDB(db)
	})

	Describe("dai", func() {
		It("writes a row", func() {
			daiMetadata := shared.StorageValueMetadata{
				Name: vat.Dai,
				Keys: map[shared.Key]string{shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, guy AS key, dai AS value FROM maker.vat_dai`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing guy", func() {
			malformedDaiMetadata := shared.StorageValueMetadata{
				Name: vat.Dai,
				Keys: nil,
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedDaiMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Guy}))
		})
	})

	Describe("gem", func() {
		It("writes row", func() {
			gemMetadata := shared.StorageValueMetadata{
				Name: vat.Gem,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk, shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedGemMetadata := shared.StorageValueMetadata{
				Name: vat.Gem,
				Keys: map[shared.Key]string{shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedGemMetadata := shared.StorageValueMetadata{
				Name: vat.Gem,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Guy}))
		})
	})

	Describe("ilk Art", func() {
		It("writes row", func() {
			ilkArtMetadata := shared.StorageValueMetadata{
				Name: vat.IlkArt,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkArtMetadata := shared.StorageValueMetadata{
				Name: vat.IlkArt,
				Keys: nil,
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})
	})

	Describe("ilk Ink", func() {
		It("writes row", func() {
			ilkInkMetadata := shared.StorageValueMetadata{
				Name: vat.IlkInk,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, ink AS value FROM maker.vat_ilk_ink`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkInkMetadata := shared.StorageValueMetadata{
				Name: vat.IlkInk,
				Keys: nil,
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})
	})

	Describe("ilk rate", func() {
		It("writes row", func() {
			ilkRateMetadata := shared.StorageValueMetadata{
				Name: vat.IlkRate,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkRateMetadata := shared.StorageValueMetadata{
				Name: vat.IlkRate,
				Keys: nil,
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkRateMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})
	})

	Describe("ilk take", func() {
		It("writes row", func() {
			ilkTakeMetadata := shared.StorageValueMetadata{
				Name: vat.IlkTake,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkTakeMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, take AS value FROM maker.vat_ilk_take`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkTakeMetadata := shared.StorageValueMetadata{
				Name: vat.IlkTake,
				Keys: nil,
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkTakeMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})
	})

	Describe("sin", func() {
		It("writes a row", func() {
			sinMetadata := shared.StorageValueMetadata{
				Name: vat.Sin,
				Keys: map[shared.Key]string{shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, guy AS key, sin AS value FROM maker.vat_sin`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing guy", func() {
			malformedSinMetadata := shared.StorageValueMetadata{
				Name: vat.Sin,
				Keys: nil,
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedSinMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Guy}))
		})
	})

	Describe("urn art", func() {
		It("writes row", func() {
			urnArtMetadata := shared.StorageValueMetadata{
				Name: vat.UrnArt,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk, shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key_one, urn AS key_two, art AS value FROM maker.vat_urn_art`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnArtMetadata := shared.StorageValueMetadata{
				Name: vat.UrnArt,
				Keys: map[shared.Key]string{shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnArtMetadata := shared.StorageValueMetadata{
				Name: vat.UrnArt,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Guy}))
		})
	})

	Describe("urn ink", func() {
		It("writes row", func() {
			urnInkMetadata := shared.StorageValueMetadata{
				Name: vat.UrnInk,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk, shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key_one, urn AS key_two, ink AS value FROM maker.vat_urn_ink`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared2.GetOrCreateIlk(fakeIlk, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeGuy, fakeUint256)
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnInkMetadata := shared.StorageValueMetadata{
				Name: vat.UrnInk,
				Keys: map[shared.Key]string{shared.Guy: fakeGuy},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnInkMetadata := shared.StorageValueMetadata{
				Name: vat.UrnInk,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Guy}))
		})
	})

	It("persists vat debt", func() {
		debtMetadata := shared.StorageValueMetadata{
			Name: vat.VatDebt,
			Keys: nil,
			Type: shared.Uint256,
		}

		err := repo.Create(fakeBlockNumber, fakeBlockHash, debtMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, debt AS value FROM maker.vat_debt`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("persists vat vice", func() {
		viceMetadata := shared.StorageValueMetadata{
			Name: vat.VatVice,
			Keys: nil,
			Type: shared.Uint256,
		}

		err := repo.Create(fakeBlockNumber, fakeBlockHash, viceMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})
})
