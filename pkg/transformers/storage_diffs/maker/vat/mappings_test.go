package vat_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/vat"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"math/big"
)

var _ = Describe("Vat storage mappings", func() {
	var (
		fakeIlk           = "fakeIlk"
		fakeGuy           = "fakeGuy"
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          vat.VatMappings
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = vat.VatMappings{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(mappings.Lookup(vat.DebtKey)).To(Equal(vat.DebtMetadata))
			Expect(mappings.Lookup(vat.ViceKey)).To(Equal(vat.ViceMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetDaiKeysCalled).To(BeTrue())
			Expect(storageRepository.GetGemKeysCalled).To(BeTrue())
			Expect(storageRepository.GetIlksCalled).To(BeTrue())
			Expect(storageRepository.GetSinKeysCalled).To(BeTrue())
			Expect(storageRepository.GetUrnsCalled).To(BeTrue())
		})

		It("returns error if dai keys lookup fails", func() {
			storageRepository.GetDaiKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if gem keys lookup fails", func() {
			storageRepository.GetGemKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if sin keys lookup fails", func() {
			storageRepository.GetSinKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if urns lookup fails", func() {
			storageRepository.GetUrnsError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("ilk", func() {
			It("returns value metadata for ilk take", func() {
				storageRepository.Ilks = []string{fakeIlk}
				ilkTakeKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.IlksMappingIndex)))
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.IlkTake,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(ilkTakeKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk rate", func() {
				storageRepository.Ilks = []string{fakeIlk}
				ilkTakeBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.IlksMappingIndex))
				ilkTakeAsInt := big.NewInt(0).SetBytes(ilkTakeBytes)
				incrementedIlkTake := big.NewInt(0).Add(ilkTakeAsInt, big.NewInt(1))
				ilkRateKey := common.BytesToHash(incrementedIlkTake.Bytes())
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.IlkRate,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(ilkRateKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk Ink", func() {
				storageRepository.Ilks = []string{fakeIlk}
				ilkTakeBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.IlksMappingIndex))
				ilkTakeAsInt := big.NewInt(0).SetBytes(ilkTakeBytes)
				doubleIncrementedIlkTake := big.NewInt(0).Add(ilkTakeAsInt, big.NewInt(2))
				ilkInkKey := common.BytesToHash(doubleIncrementedIlkTake.Bytes())
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.IlkInk,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(ilkInkKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk Art", func() {
				storageRepository.Ilks = []string{fakeIlk}
				ilkTakeBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.IlksMappingIndex))
				ilkTakeAsInt := big.NewInt(0).SetBytes(ilkTakeBytes)
				tripleIncrementedIlkTake := big.NewInt(0).Add(ilkTakeAsInt, big.NewInt(3))
				ilkArtKey := common.BytesToHash(tripleIncrementedIlkTake.Bytes())
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.IlkArt,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(ilkArtKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("urn", func() {
			It("returns value metadata for urn ink", func() {
				storageRepository.Urns = []maker.Urn{{Ilk: fakeIlk, Guy: fakeGuy}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.UrnsMappingIndex))
				encodedSecondaryMapIndex := crypto.Keccak256(common.FromHex(fakeGuy), encodedPrimaryMapIndex)
				urnInkKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.UrnInk,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk, shared.Guy: fakeGuy},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(urnInkKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for urn art", func() {
				storageRepository.Urns = []maker.Urn{{Ilk: fakeIlk, Guy: fakeGuy}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.UrnsMappingIndex))
				urnInkAsInt := big.NewInt(0).SetBytes(crypto.Keccak256(common.FromHex(fakeGuy), encodedPrimaryMapIndex))
				incrementedUrnInk := big.NewInt(0).Add(urnInkAsInt, big.NewInt(1))
				urnArtKey := common.BytesToHash(incrementedUrnInk.Bytes())
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.UrnArt,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk, shared.Guy: fakeGuy},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(urnArtKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("gem", func() {
			It("returns value metadata for gem", func() {
				storageRepository.GemKeys = []maker.Urn{{Ilk: fakeIlk, Guy: fakeGuy}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.GemsMappingIndex))
				encodedSecondaryMapIndex := crypto.Keccak256(common.FromHex(fakeGuy), encodedPrimaryMapIndex)
				gemKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.Gem,
					Keys: map[shared.Key]string{shared.Ilk: fakeIlk, shared.Guy: fakeGuy},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(gemKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("dai", func() {
			It("returns value metadata for dai", func() {
				storageRepository.DaiKeys = []string{fakeGuy}
				daiKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeGuy + vat.DaiMappingIndex)))
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.Dai,
					Keys: map[shared.Key]string{shared.Guy: fakeGuy},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(daiKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("when sin key exists in the db", func() {
			It("returns value metadata for sin", func() {
				storageRepository.SinKeys = []string{fakeGuy}
				sinKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeGuy + vat.SinMappingIndex)))
				expectedMetadata := shared.StorageValueMetadata{
					Name: vat.Sin,
					Keys: map[shared.Key]string{shared.Guy: fakeGuy},
					Type: shared.Uint256,
				}

				Expect(mappings.Lookup(sinKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
