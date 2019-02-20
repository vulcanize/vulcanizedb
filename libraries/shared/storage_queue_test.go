package shared_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	shared2 "github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage queue", func() {
	It("adds a storage row to the db", func() {
		row := shared.StorageDiffRow{
			Contract:     common.HexToAddress("0x123456"),
			BlockHash:    common.HexToHash("0x678901"),
			BlockHeight:  987,
			StorageKey:   common.HexToHash("0x654321"),
			StorageValue: common.HexToHash("0x198765"),
		}
		db := test_config.NewTestDB(test_config.NewTestNode())
		queue := shared2.NewStorageQueue(db)

		addErr := queue.Add(row)

		Expect(addErr).NotTo(HaveOccurred())
		var result shared.StorageDiffRow
		getErr := db.Get(&result, `SELECT contract, block_hash, block_height, storage_key, storage_value FROM public.queued_storage`)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(result).To(Equal(row))
	})
})
