package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
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

	It("returns metadata for a packed storage slot variables", func() {
		metadataName := "fake_name"
		metadataKeys := map[utils.Key]string{"key": "value"}
		metadataType := utils.Uint256
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
})
