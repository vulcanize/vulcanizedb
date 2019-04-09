package storage_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage queue", func() {
	It("adds a storage row to the db", func() {
		row := utils.StorageDiffRow{
			Contract:     common.HexToAddress("0x123456"),
			BlockHash:    common.HexToHash("0x678901"),
			BlockHeight:  987,
			StorageKey:   common.HexToHash("0x654321"),
			StorageValue: common.HexToHash("0x198765"),
		}
		db := test_config.NewTestDB(test_config.NewTestNode())
		queue := storage.NewStorageQueue(db)

		addErr := queue.Add(row)

		Expect(addErr).NotTo(HaveOccurred())
		var result utils.StorageDiffRow
		getErr := db.Get(&result, `SELECT contract, block_hash, block_height, storage_key, storage_value FROM public.queued_storage`)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(result).To(Equal(row))
	})
})
