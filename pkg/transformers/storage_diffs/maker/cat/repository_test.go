package cat_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/cat"
	. "github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Cat storage repository", func() {
	var (
		db              *postgres.DB
		repo            cat.CatStorageRepository
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeAddress     = "0x12345"
		fakeIlk         = "fake_ilk"
		fakeUint256     = "12345"
		fakeBytes32     = "fake_bytes32"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = cat.CatStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Variable", func() {
		var result VariableRes

		Describe("NFlip", func() {
			It("writes a row", func() {
				nFlipMetadata := shared.GetStorageValueMetadata(cat.NFlip, nil, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, nFlipMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, nflip AS value FROM maker.cat_nflip`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
			})
		})

		Describe("Live", func() {
			It("writes a row", func() {
				liveMetadata := shared.GetStorageValueMetadata(cat.Live, nil, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, liveMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.cat_live`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
			})
		})

		Describe("Vat", func() {
			It("writes a row", func() {
				vatMetadata := shared.GetStorageValueMetadata(cat.Vat, nil, shared.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, vatMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.cat_vat`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
			})
		})

		Describe("Pit", func() {
			It("writes a row", func() {
				pitMetadata := shared.GetStorageValueMetadata(cat.Pit, nil, shared.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, pitMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, pit AS value FROM maker.cat_pit`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
			})
		})

		Describe("Vow", func() {
			It("writes a row", func() {
				vowMetadata := shared.GetStorageValueMetadata(cat.Vow, nil, shared.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, vowMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, vow AS value FROM maker.cat_vow`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
			})
		})
	})

	Describe("Ilk", func() {
		var result MappingRes

		Describe("Flip", func() {
			It("writes a row", func() {
				ilkFlipMetadata := shared.GetStorageValueMetadata(cat.IlkFlip, map[shared.Key]string{shared.Ilk: fakeIlk}, shared.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkFlipMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, flip AS value FROM maker.cat_ilk_flip`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeIlk, fakeAddress)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkFlipMetadata := shared.GetStorageValueMetadata(cat.IlkFlip, map[shared.Key]string{}, shared.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkFlipMetadata, fakeAddress)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
			})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := shared.GetStorageValueMetadata(cat.IlkChop, map[shared.Key]string{shared.Ilk: fakeIlk}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkChopMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, chop AS value FROM maker.cat_ilk_chop`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeIlk, fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkChopMetadata := shared.GetStorageValueMetadata(cat.IlkChop, map[shared.Key]string{}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkChopMetadata, fakeAddress)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
			})
		})

		Describe("Lump", func() {
			It("writes a row", func() {
				ilkLumpMetadata := shared.GetStorageValueMetadata(cat.IlkLump, map[shared.Key]string{shared.Ilk: fakeIlk}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLumpMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, lump AS value FROM maker.cat_ilk_lump`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeIlk, fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkLumpMetadata := shared.GetStorageValueMetadata(cat.IlkLump, map[shared.Key]string{}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLumpMetadata, fakeAddress)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Ilk}))
			})
		})
	})

	Describe("Flip", func() {
		var result MappingRes

		Describe("FlipIlk", func() {
			It("writes a row", func() {
				flipIlkMetadata := shared.GetStorageValueMetadata(cat.FlipIlk, map[shared.Key]string{shared.Flip: fakeUint256}, shared.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipIlkMetadata, fakeBytes32)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, nflip AS key, ilk AS value FROM maker.cat_flip_ilk`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeBytes32)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipIlkMetadata := shared.GetStorageValueMetadata(cat.FlipIlk, map[shared.Key]string{}, shared.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipIlkMetadata, fakeBytes32)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Flip}))
			})
		})

		Describe("FlipUrn", func() {
			It("writes a row", func() {
				flipUrnMetadata := shared.GetStorageValueMetadata(cat.FlipUrn, map[shared.Key]string{shared.Flip: fakeUint256}, shared.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipUrnMetadata, fakeBytes32)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, nflip AS key, urn AS value FROM maker.cat_flip_urn`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeBytes32)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipUrnMetadata := shared.GetStorageValueMetadata(cat.FlipUrn, map[shared.Key]string{}, shared.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipUrnMetadata, fakeBytes32)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Flip}))
			})
		})

		Describe("FlipInk", func() {
			It("writes a row", func() {
				flipInkMetadata := shared.GetStorageValueMetadata(cat.FlipInk, map[shared.Key]string{shared.Flip: fakeUint256}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipInkMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, nflip AS key, ink AS value FROM maker.cat_flip_ink`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeUint256)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipInkMetadata := shared.GetStorageValueMetadata(cat.FlipInk, map[shared.Key]string{}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipInkMetadata, fakeUint256)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Flip}))
			})
		})

		Describe("FlipTab", func() {
			It("writes a row", func() {
				flipTabMetadata := shared.GetStorageValueMetadata(cat.FlipTab, map[shared.Key]string{shared.Flip: fakeUint256}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipTabMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, nflip AS key, tab AS value FROM maker.cat_flip_tab`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeUint256)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipTabMetadata := shared.GetStorageValueMetadata(cat.FlipTab, map[shared.Key]string{}, shared.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipTabMetadata, fakeUint256)
				Expect(err).To(MatchError(shared.ErrMetadataMalformed{MissingData: shared.Flip}))
			})
		})
	})
})
