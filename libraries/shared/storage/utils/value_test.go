package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

var _ = Describe("Storage value metadata getter", func() {
	It("returns a storage value metadata instance with corresponding fields assigned", func() {
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
})
