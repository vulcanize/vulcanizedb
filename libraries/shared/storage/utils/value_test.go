package utils_test

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage value metadata getter", func() {
	It("returns storage value metadata for a single storage variable", func() {
		metadataName := "fake_name"
		metadataKeys := map[utils.Key]string{"key": "value"}
		metadataType := utils.Uint256

		expectedMetadata := utils.StorageValueMetadata{
			Name: metadataName,
			Keys: metadataKeys,
			Type: metadataType,
		}
		Expect(utils.GetStorageValueMetadata(metadataName, metadataKeys, metadataType)).To(Equal(expectedMetadata))
	})

	Describe("metadata for a packed storaged slot", func() {
		It("returns metadata for multiple storage variables", func() {
			metadataName := "fake_name"
			metadataKeys := map[utils.Key]string{"key": "value"}
			metadataType := utils.PackedSlot
			metadataPackedNames := map[int]string{0: "name"}
			metadataPackedTypes := map[int]utils.ValueType{0: utils.Uint48}

			expectedMetadata := utils.StorageValueMetadata{
				Name:        metadataName,
				Keys:        metadataKeys,
				Type:        metadataType,
				PackedTypes: metadataPackedTypes,
				PackedNames: metadataPackedNames,
			}
			Expect(utils.GetStorageValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, metadataPackedNames, metadataPackedTypes)).To(Equal(expectedMetadata))
		})

		It("panics if PackedTypes are nil when the type is PackedSlot", func() {
			metadataName := "fake_name"
			metadataKeys := map[utils.Key]string{"key": "value"}
			metadataType := utils.PackedSlot
			metadataPackedNames := map[int]string{0: "name"}

			getMetadata := func() {
				utils.GetStorageValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, metadataPackedNames, nil)
			}
			Expect(getMetadata).To(Panic())
		})

		It("panics if PackedNames are nil when the type is PackedSlot", func() {
			metadataName := "fake_name"
			metadataKeys := map[utils.Key]string{"key": "value"}
			metadataType := utils.PackedSlot
			metadataPackedTypes := map[int]utils.ValueType{0: utils.Uint48}

			getMetadata := func() {
				utils.GetStorageValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, nil, metadataPackedTypes)
			}
			Expect(getMetadata).To(Panic())
		})

		It("panics if valueType is not PackedSlot if PackedNames is populated", func() {
			metadataName := "fake_name"
			metadataKeys := map[utils.Key]string{"key": "value"}
			metadataType := utils.Uint48
			metadataPackedNames := map[int]string{0: "name"}
			metadataPackedTypes := map[int]utils.ValueType{0: utils.Uint48}

			getMetadata := func() {
				utils.GetStorageValueMetadataForPackedSlot(metadataName, metadataKeys, metadataType, metadataPackedNames, metadataPackedTypes)
			}
			Expect(getMetadata).To(Panic())
		})
	})
})
