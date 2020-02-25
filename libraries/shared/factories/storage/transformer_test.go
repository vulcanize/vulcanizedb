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

package storage_test

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/mocks"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage transformer", func() {
	var (
		storageKeysLookup *mocks.MockStorageKeysLookup
		repository        *mocks.MockStorageRepository
		t                 storage.Transformer
	)

	BeforeEach(func() {
		storageKeysLookup = &mocks.MockStorageKeysLookup{}
		repository = &mocks.MockStorageRepository{}
		t = storage.Transformer{
			Address:           common.Address{},
			StorageKeysLookup: storageKeysLookup,
			Repository:        repository,
		}
	})

	It("returns the keccaked contract address being watched", func() {
		fakeAddress := fakes.FakeAddress
		keccakOfAddress := types.HexToKeccak256Hash(fakeAddress.Hex())
		t.Address = fakeAddress

		Expect(t.KeccakContractAddress()).To(Equal(keccakOfAddress))
	})

	It("returns the contract address being watched", func() {
		fakeAddress := fakes.FakeAddress
		t.Address = fakeAddress

		Expect(t.GetContractAddress()).To(Equal(fakeAddress))
	})

	It("looks up metadata for storage key", func() {
		t.Execute(types.PersistedDiff{})

		Expect(storageKeysLookup.LookupCalled).To(BeTrue())
	})

	It("returns error if lookup fails", func() {
		storageKeysLookup.LookupErr = fakes.FakeError

		err := t.Execute(types.PersistedDiff{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("creates storage row with decoded data", func() {
		fakeMetadata := types.ValueMetadata{Type: types.Address}
		storageKeysLookup.Metadata = fakeMetadata
		rawValue := common.HexToAddress("0x12345")
		fakeHeaderID := rand.Int63()
		fakeBlockNumber := rand.Int()
		fakeBlockHash := fakes.RandomString(64)
		fakeRow := types.PersistedDiff{
			ID:       rand.Int63(),
			HeaderID: fakeHeaderID,
			RawDiff: types.RawDiff{
				HashedAddress: common.Hash{},
				BlockHash:     common.HexToHash(fakeBlockHash),
				BlockHeight:   fakeBlockNumber,
				StorageKey:    common.Hash{},
				StorageValue:  rawValue.Hash(),
			},
		}

		err := t.Execute(fakeRow)

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeaderID))
		Expect(repository.PassedDiffID).To(Equal(fakeRow.ID))
		Expect(repository.PassedMetadata).To(Equal(fakeMetadata))
		Expect(repository.PassedValue.(string)).To(Equal(rawValue.Hex()))
	})

	It("returns error if creating row fails", func() {
		rawValue := common.HexToAddress("0x12345")
		fakeMetadata := types.ValueMetadata{Type: types.Address}
		storageKeysLookup.Metadata = fakeMetadata
		repository.CreateErr = fakes.FakeError
		diff := types.PersistedDiff{RawDiff: types.RawDiff{StorageValue: rawValue.Hash()}}

		err := t.Execute(diff)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	Describe("when a storage row contains more than one item packed in storage", func() {
		var (
			rawValue     = common.HexToAddress("000000000000000000000000000000000000000000000002a300000000002a30")
			fakeHeaderID = rand.Int63()
			packedTypes  = make(map[int]types.ValueType)
		)
		packedTypes[0] = types.Uint48
		packedTypes[1] = types.Uint48

		var fakeMetadata = types.ValueMetadata{
			Name:        "",
			Keys:        nil,
			Type:        types.PackedSlot,
			PackedTypes: packedTypes,
		}

		It("passes the decoded data items to the repository", func() {
			storageKeysLookup.Metadata = fakeMetadata
			fakeBlockNumber := rand.Int()
			fakeBlockHash := fakes.RandomString(64)
			fakeRow := types.PersistedDiff{
				ID:       rand.Int63(),
				HeaderID: fakeHeaderID,
				RawDiff: types.RawDiff{
					HashedAddress: common.Hash{},
					BlockHash:     common.HexToHash(fakeBlockHash),
					BlockHeight:   fakeBlockNumber,
					StorageKey:    common.Hash{},
					StorageValue:  rawValue.Hash(),
				},
			}

			err := t.Execute(fakeRow)

			Expect(err).NotTo(HaveOccurred())
			Expect(repository.PassedHeaderID).To(Equal(fakeHeaderID))
			Expect(repository.PassedDiffID).To(Equal(fakeRow.ID))
			Expect(repository.PassedMetadata).To(Equal(fakeMetadata))
			expectedPassedValue := make(map[int]string)
			expectedPassedValue[0] = "10800"
			expectedPassedValue[1] = "172800"
			Expect(repository.PassedValue.(map[int]string)).To(Equal(expectedPassedValue))
		})

		It("returns error if creating a row fails", func() {
			storageKeysLookup.Metadata = fakeMetadata
			repository.CreateErr = fakes.FakeError
			diff := types.PersistedDiff{RawDiff: types.RawDiff{StorageValue: rawValue.Hash()}}

			err := t.Execute(diff)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})
})
