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

package cold_import_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth/cold_import"
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

	It("only populates missing blocks", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		nodeId := "node_id"
		startingBlockNumber := int64(120)
		missingBlockNumber := int64(123)
		endingBlockNumber := int64(125)
		fakeHash := []byte{1, 2, 3, 4, 5}
		mockBlockRepository.SetMissingBlockNumbersReturnArray([]int64{missingBlockNumber})
		mockEthereumDatabase.SetReturnHash(fakeHash)
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		importer := cold_import.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		importer.Execute(startingBlockNumber, endingBlockNumber, nodeId)

		mockBlockRepository.AssertMissingBlockNumbersCalledWith(startingBlockNumber, endingBlockNumber, nodeId)
		mockEthereumDatabase.AssertGetBlockHashCalledWith(missingBlockNumber)
		mockEthereumDatabase.AssertGetBlockCalledWith(fakeHash, missingBlockNumber)
	})

	It("fetches missing blocks from level db and persists them to pg", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		blockNumber := int64(123)
		fakeHash := []byte{1, 2, 3, 4, 5}
		mockBlockRepository.SetMissingBlockNumbersReturnArray([]int64{blockNumber})
		mockEthereumDatabase.SetReturnHash(fakeHash)
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		importer := cold_import.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		importer.Execute(blockNumber, blockNumber, "node_id")

		mockEthereumDatabase.AssertGetBlockHashCalledWith(blockNumber)
		mockEthereumDatabase.AssertGetBlockCalledWith(fakeHash, blockNumber)
		mockTransactionConverter.AssertConvertTransactionsToCoreCalledWith(fakeGethBlock)
		convertedBlock, err := blockConverter.ToCoreBlock(fakeGethBlock)
		Expect(err).NotTo(HaveOccurred())
		mockBlockRepository.AssertCreateOrUpdateBlockCalledWith(convertedBlock)
	})

	It("sets is_final status on populated blocks", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		startingBlockNumber := int64(120)
		endingBlockNumber := int64(125)
		fakeHash := []byte{1, 2, 3, 4, 5}
		mockBlockRepository.SetMissingBlockNumbersReturnArray([]int64{startingBlockNumber})
		mockEthereumDatabase.SetReturnHash(fakeHash)
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		importer := cold_import.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		importer.Execute(startingBlockNumber, endingBlockNumber, "node_id")

		mockBlockRepository.AssertSetBlockStatusCalledWith(endingBlockNumber)
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
		mockBlockRepository.SetMissingBlockNumbersReturnArray([]int64{blockNumber})
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		mockEthereumDatabase.SetReturnReceipts(fakeReceipts)
		importer := cold_import.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		importer.Execute(blockNumber, blockNumber, "node_id")

		expectedReceipts := vulcCommon.ToCoreReceipts(fakeReceipts)
		mockReceiptRepository.AssertCreateReceiptsAndLogsCalledWith(blockId, expectedReceipts)
	})

	It("does not fetch receipts if block already exists", func() {
		mockEthereumDatabase := fakes.NewMockEthereumDatabase()
		mockBlockRepository := fakes.NewMockBlockRepository()
		mockReceiptRepository := fakes.NewMockReceiptRepository()
		mockTransactionConverter := fakes.NewMockTransactionConverter()
		blockConverter := vulcCommon.NewBlockConverter(mockTransactionConverter)

		blockNumber := int64(123)
		mockBlockRepository.SetMissingBlockNumbersReturnArray([]int64{})
		mockEthereumDatabase.SetReturnBlock(fakeGethBlock)
		mockBlockRepository.SetCreateOrUpdateBlockReturnVals(0, repositories.ErrBlockExists)
		importer := cold_import.NewColdImporter(mockEthereumDatabase, mockBlockRepository, mockReceiptRepository, blockConverter)

		err := importer.Execute(blockNumber, blockNumber, "node_id")

		Expect(err).NotTo(HaveOccurred())
		mockReceiptRepository.AssertCreateReceiptsAndLogsNotCalled()
	})
})
