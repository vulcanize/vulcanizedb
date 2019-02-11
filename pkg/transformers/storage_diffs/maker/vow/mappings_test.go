package vow_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker/vow"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

var _ = Describe("Vow storage mappings", func() {
	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}

			mappings := vow.VowMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(vow.VatKey)).To(Equal(vow.VatMetadata))
			Expect(mappings.Lookup(vow.CowKey)).To(Equal(vow.CowMetadata))
			Expect(mappings.Lookup(vow.RowKey)).To(Equal(vow.RowMetadata))
			Expect(mappings.Lookup(vow.SinKey)).To(Equal(vow.SinMetadata))
			Expect(mappings.Lookup(vow.AshKey)).To(Equal(vow.AshMetadata))
			Expect(mappings.Lookup(vow.WaitKey)).To(Equal(vow.WaitMetadata))
			Expect(mappings.Lookup(vow.SumpKey)).To(Equal(vow.SumpMetadata))
			Expect(mappings.Lookup(vow.BumpKey)).To(Equal(vow.BumpMetadata))
			Expect(mappings.Lookup(vow.HumpKey)).To(Equal(vow.HumpMetadata))
		})

		It("returns error if key does not exist", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}

			mappings := vow.VowMappings{StorageRepository: storageRepository}
			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})
})
