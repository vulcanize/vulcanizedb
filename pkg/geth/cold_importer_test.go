package geth_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	vulcCommon "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
)

var _ = Describe("Geth cold importer", func() {
	var fakeGethBlock *types.Block

	BeforeEach(func() {
		header := &types.Header{}
		transactions := []*types.Transaction{}
		uncles := []*types.Header{}
		receipts := []*types.Receipt{}
		fakeGethBlock = types.NewBlock(header, transactions, uncles, receipts)
	})

	It("fetches blocks from level db and persists them to pg", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		blockNumber := int64(123)
		fakeHash := []byte{1, 2, 3, 4, 5}
		mockEthereumDatabase.SetReturnHash(fakeHash)
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		importer := geth.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		importer.Execute(blockNumber, blockNumber)

		mockEthereumDatabase.AssertGetBlockHashCalledWith(blockNumber)
		mockEthereumDatabase.AssertGetBlockCalledWith(fakeHash, blockNumber)
		mockTransactionConverter.AssertConvertTransactionsToCoreCalledWith(fakeGethBlock)
		convertedBlock, err := blockConverter.ToCoreBlock(fakeGethBlock)
		Expect(err).NotTo(HaveOccurred())
		mockBlockRepository.AssertCreateOrUpdateBlockCalledWith(convertedBlock)
	})

	It("fetches receipts from level db and persists them to pg", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		blockNumber := int64(123)
		blockId := int64(999)
		mockBlockRepository.SetCreateOrUpdateBlockReturnVals(blockId, nil)
		fakeReceipts := types.Receipts{{}}
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		mockEthereumDatabase.SetReturnReceipts(fakeReceipts)
		importer := geth.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		importer.Execute(blockNumber, blockNumber)

		expectedReceipts := vulcCommon.ToCoreReceipts(fakeReceipts)
		mockReceiptRepository.AssertCreateReceiptsAndLogsCalledWith(blockId, expectedReceipts)
	})

	It("does not fetch receipts if block already exists", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		mockBlockRepository.SetCreateOrUpdateBlockReturnVals(0, repositories.ErrBlockExists)
		importer := geth.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		err := importer.Execute(1, 1)

		Expect(err).NotTo(HaveOccurred())
		mockReceiptRepository.AssertCreateReceiptsAndLogsNotCalled()
	})
})
