package pit_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/pit"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"math/big"
)

var _ = Describe("Pit storage mappings", func() {
	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := pit.PitMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(pit.DripKey)).To(Equal(pit.DripMetadata))
			Expect(mappings.Lookup(pit.LineKey)).To(Equal(pit.LineMetadata))
			Expect(mappings.Lookup(pit.LiveKey)).To(Equal(pit.LiveMetadata))
			Expect(mappings.Lookup(pit.VatKey)).To(Equal(pit.VatMetadata))
		})

		It("returns error if key does not exist", func() {
			mappings := pit.PitMappings{StorageRepository: &test_helpers.MockMakerStorageRepository{}}

			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := pit.PitMappings{StorageRepository: storageRepository}

			mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
		})

		It("returns value metadata for spot when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.SetIlks([]string{fakeIlk})
			mappings := pit.PitMappings{StorageRepository: storageRepository}
			ilkSpotKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + pit.IlkSpotIndex)))
			expectedMetadata := shared.StorageValueMetadata{
				Name: pit.IlkSpot,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			Expect(mappings.Lookup(ilkSpotKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for line when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.SetIlks([]string{fakeIlk})
			mappings := pit.PitMappings{StorageRepository: storageRepository}
			ilkSpotKeyBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + pit.IlkSpotIndex))
			ilkSpotAsInt := big.NewInt(0).SetBytes(ilkSpotKeyBytes)
			incrementedIlkSpot := big.NewInt(0).Add(ilkSpotAsInt, big.NewInt(1))
			ilkLineKey := common.BytesToHash(incrementedIlkSpot.Bytes())
			expectedMetadata := shared.StorageValueMetadata{
				Name: pit.IlkLine,
				Keys: map[shared.Key]string{shared.Ilk: fakeIlk},
				Type: shared.Uint256,
			}

			Expect(mappings.Lookup(ilkLineKey)).To(Equal(expectedMetadata))
		})

		It("returns error if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := pit.PitMappings{StorageRepository: storageRepository}

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})
})
