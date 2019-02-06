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

package storage_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories/storage"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
)

var _ = Describe("Storage transformer", func() {
	var (
		mappings    *mocks.MockMappings
		repository  *mocks.MockStorageRepository
		transformer storage.Transformer
	)

	BeforeEach(func() {
		mappings = &mocks.MockMappings{}
		repository = &mocks.MockStorageRepository{}
		transformer = storage.Transformer{
			Address:    common.Address{},
			Mappings:   mappings,
			Repository: repository,
		}
	})

	It("returns the contract address being watched", func() {
		fakeAddress := common.HexToAddress("0x12345")
		transformer.Address = fakeAddress

		Expect(transformer.ContractAddress()).To(Equal(fakeAddress))
	})

	It("looks up metadata for storage key", func() {
		transformer.Execute(shared.StorageDiffRow{})

		Expect(mappings.LookupCalled).To(BeTrue())
	})

	It("returns error if lookup fails", func() {
		mappings.LookupErr = fakes.FakeError

		err := transformer.Execute(shared.StorageDiffRow{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("creates storage row with decoded data", func() {
		fakeMetadata := shared.StorageValueMetadata{Type: shared.Address}
		mappings.Metadata = fakeMetadata
		rawValue := common.HexToAddress("0x12345")
		fakeBlockNumber := 123
		fakeBlockHash := "0x67890"
		fakeRow := shared.StorageDiffRow{
			Contract:     common.Address{},
			BlockHash:    common.HexToHash(fakeBlockHash),
			BlockHeight:  fakeBlockNumber,
			StorageKey:   common.Hash{},
			StorageValue: rawValue.Hash(),
		}

		err := transformer.Execute(fakeRow)

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedBlockNumber).To(Equal(fakeBlockNumber))
		Expect(repository.PassedBlockHash).To(Equal(common.HexToHash(fakeBlockHash).Hex()))
		Expect(repository.PassedMetadata).To(Equal(fakeMetadata))
		Expect(repository.PassedValue.(string)).To(Equal(rawValue.Hex()))
	})

	It("returns error if creating row fails", func() {
		rawValue := common.HexToAddress("0x12345")
		fakeMetadata := shared.StorageValueMetadata{Type: shared.Address}
		mappings.Metadata = fakeMetadata
		repository.CreateErr = fakes.FakeError

		err := transformer.Execute(shared.StorageDiffRow{StorageValue: rawValue.Hash()})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
