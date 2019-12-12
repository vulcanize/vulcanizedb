package storage_test

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage value metadata getter", func() {
	It("returns storage value metadata for a single storage variable", func() {
		metadataName := "fake_name"
		metadataKeys := map[storage.Key]string{"key": "value"}
		metadataType := storage.Uint256

		expectedMetadata := storage.ValueMetadata{
			Name: metadataName,
			Keys: metadataKeys,
			Type: metadataType,
		}
		Expect(storage.GetValueMetadata(metadataName, metadataKeys, metadataType)).To(Equal(expectedMetadata))
	})

	Describe("metadata for a packed storaged slot", func() {
		It("returns metadata for multiple storage variables", func() {
			metadataName := "fake_name"
			metadataKeys := map[storage.Key]string{"key": "value"}
			metadataType := storage.PackedSlot
			metadataPackedNames := map[int]string{0: "name"}
			metadataPackedTypes := map[int]storage.ValueType{0: storage.Uint48}

			expectedMetadata := storage.ValueMetadata{
				Name:        metadataName,
				Keys:        metadataKeys,
				Type:        metadataType,
				PackedTypes: metadataPackedTypes,
				PackedNames: metadataPackedNames,
			}
			Expect(storage.GetValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, metadataPackedNames, metadataPackedTypes)).To(Equal(expectedMetadata))
		})

		It("panics if PackedTypes are nil when the type is PackedSlot", func() {
			metadataName := "fake_name"
			metadataKeys := map[storage.Key]string{"key": "value"}
			metadataType := storage.PackedSlot
			metadataPackedNames := map[int]string{0: "name"}

			getMetadata := func() {
				storage.GetValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, metadataPackedNames, nil)
			}
			Expect(getMetadata).To(Panic())
		})

		It("panics if PackedNames are nil when the type is PackedSlot", func() {
			metadataName := "fake_name"
			metadataKeys := map[storage.Key]string{"key": "value"}
			metadataType := storage.PackedSlot
			metadataPackedTypes := map[int]storage.ValueType{0: storage.Uint48}

			getMetadata := func() {
				storage.GetValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, nil, metadataPackedTypes)
			}
			Expect(getMetadata).To(Panic())
		})

		It("panics if valueType is not PackedSlot if PackedNames is populated", func() {
			metadataName := "fake_name"
			metadataKeys := map[storage.Key]string{"key": "value"}
			metadataType := storage.Uint48
			metadataPackedNames := map[int]string{0: "name"}
			metadataPackedTypes := map[int]storage.ValueType{0: storage.Uint48}

			getMetadata := func() {
				storage.GetValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, metadataPackedNames, metadataPackedTypes)
			}
			Expect(getMetadata).To(Panic())
		})
	})
})
