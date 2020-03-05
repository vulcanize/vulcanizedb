// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package repositories_test

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Block header repository", func() {
	var (
		db     = test_config.NewTestDB(test_config.NewTestNode())
		repo   repositories.HeaderRepository
		header core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = repositories.NewHeaderRepository(db)
		header = fakes.GetFakeHeader(rand.Int63n(50000000))
	})

	Describe("creating or updating a header", func() {
		BeforeEach(func() {
			_, createErr := repo.CreateOrUpdateHeader(header)
			Expect(createErr).NotTo(HaveOccurred())
		})

		It("adds a header", func() {
			var dbHeader core.Header
			readErr := db.Get(&dbHeader, `SELECT block_number, hash, raw, block_timestamp FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(dbHeader.BlockNumber).To(Equal(header.BlockNumber))
			Expect(dbHeader.Hash).To(Equal(header.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(header.Raw))
			Expect(dbHeader.Timestamp).To(Equal(header.Timestamp))
		})

		It("adds node data to header", func() {
			var ethNodeId int64
			readErr := db.Get(&ethNodeId, `SELECT eth_node_id FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(ethNodeId).To(Equal(db.NodeID))
		})

		It("does not duplicate headers", func() {
			_, createTwoErr := repo.CreateOrUpdateHeader(header)
			Expect(createTwoErr).NotTo(HaveOccurred())

			var count int
			readErr := db.Get(&count, `SELECT COUNT(*) FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("replaces header if hash is different and block number is within 15 of the max block number in the db", func() {
			headerTwo := fakes.GetFakeHeader(header.BlockNumber)

			_, createTwoErr := repo.CreateOrUpdateHeader(headerTwo)

			Expect(createTwoErr).NotTo(HaveOccurred())
			var dbHeaderHash string
			readErr := db.Get(&dbHeaderHash, `SELECT hash FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(dbHeaderHash).To(Equal(headerTwo.Hash))
		})

		It("does not replace header if block number if greater than 15 back from the max block number in the db", func() {
			chainHeadHeader := fakes.GetFakeHeader(header.BlockNumber + 15)
			_, createHeadErr := repo.CreateOrUpdateHeader(chainHeadHeader)
			Expect(createHeadErr).NotTo(HaveOccurred())

			oldConflictingHeader := fakes.GetFakeHeader(header.BlockNumber)
			_, createConflictErr := repo.CreateOrUpdateHeader(oldConflictingHeader)
			Expect(createConflictErr).NotTo(HaveOccurred())

			var dbHeaderHash string
			readErr := db.Get(&dbHeaderHash, `SELECT hash FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(dbHeaderHash).To(Equal(header.Hash))
		})

		It("does not duplicate headers with different hashes", func() {
			headerTwo := fakes.GetFakeHeader(header.BlockNumber)

			_, createTwoErr := repo.CreateOrUpdateHeader(headerTwo)
			Expect(createTwoErr).NotTo(HaveOccurred())

			var dbHeaderHashes []string
			readErr := db.Select(&dbHeaderHashes, `SELECT hash FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(len(dbHeaderHashes)).To(Equal(1))
			Expect(dbHeaderHashes[0]).To(Equal(headerTwo.Hash))
		})

		It("replaces header if hash is different (even from different node)", func() {
			dbTwo := test_config.NewTestDB(test_config.NewTestNode())

			repoTwo := repositories.NewHeaderRepository(dbTwo)
			headerTwo := fakes.GetFakeHeader(header.BlockNumber)

			_, createTwoErr := repoTwo.CreateOrUpdateHeader(headerTwo)

			Expect(createTwoErr).NotTo(HaveOccurred())
			var dbHeaderHashes []string
			readErr := dbTwo.Select(&dbHeaderHashes, `SELECT hash FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(len(dbHeaderHashes)).To(Equal(1))
			Expect(dbHeaderHashes[0]).To(Equal(headerTwo.Hash))
		})
	})

	Describe("creating a receipt", func() {
		It("adds a receipt in a tx", func() {
			headerID, err := repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := common.HexToAddress("0x1234")
			toAddress := common.HexToAddress("0x5678")
			txHash := common.HexToHash("0x9876")
			txIndex := big.NewInt(123)
			transaction := core.TransactionModel{
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 0,
				GasPrice: 0,
				Hash:     txHash.Hex(),
				Nonce:    0,
				Raw:      []byte{},
				To:       toAddress.Hex(),
				TxIndex:  txIndex.Int64(),
				Value:    "0",
			}
			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())
			txId, txErr := repo.CreateTransactionInTx(tx, headerID, transaction)
			Expect(txErr).ToNot(HaveOccurred())

			contractAddr := common.HexToAddress("0x1234")
			stateRoot := common.HexToHash("0x5678")
			receipt := core.Receipt{
				ContractAddress:   contractAddr.Hex(),
				TxHash:            txHash.Hex(),
				GasUsed:           10,
				CumulativeGasUsed: 100,
				StateRoot:         stateRoot.Hex(),
				Rlp:               []byte{1, 2, 3},
			}

			receiptRepo := repositories.ReceiptRepository{}
			_, receiptErr := receiptRepo.CreateReceiptInTx(headerID, txId, receipt, tx)
			Expect(receiptErr).ToNot(HaveOccurred())
			commitErr := tx.Commit()
			Expect(commitErr).ToNot(HaveOccurred())

			type idModel struct {
				TransactionId     int64  `db:"transaction_id"`
				ContractAddressId int64  `db:"contract_address_id"`
				CumulativeGasUsed uint64 `db:"cumulative_gas_used"`
				GasUsed           uint64 `db:"gas_used"`
				StateRoot         string `db:"state_root"`
				Status            int
				TxHash            string `db:"tx_hash"`
				Rlp               []byte `db:"rlp"`
			}

			var addressId int64
			getAddressErr := db.Get(&addressId, `SELECT id FROM addresses WHERE address = $1`, contractAddr.Hex())
			Expect(getAddressErr).NotTo(HaveOccurred())

			var dbReceipt idModel
			getReceiptErr := db.Get(&dbReceipt,
				`SELECT transaction_id, contract_address_id, cumulative_gas_used, gas_used, state_root, status, tx_hash, rlp
				FROM public.receipts WHERE header_id = $1`, headerID)
			Expect(getReceiptErr).NotTo(HaveOccurred())

			Expect(dbReceipt.TransactionId).To(Equal(txId))
			Expect(dbReceipt.TxHash).To(Equal(txHash.Hex()))
			Expect(dbReceipt.ContractAddressId).To(Equal(addressId))
			Expect(dbReceipt.CumulativeGasUsed).To(Equal(uint64(100)))
			Expect(dbReceipt.GasUsed).To(Equal(uint64(10)))
			Expect(dbReceipt.StateRoot).To(Equal(stateRoot.Hex()))
			Expect(dbReceipt.Status).To(Equal(0))
			Expect(dbReceipt.Rlp).To(Equal([]byte{1, 2, 3}))
		})
	})

	Describe("creating a transaction", func() {
		var (
			headerID     int64
			transactions []core.TransactionModel
		)

		BeforeEach(func() {
			var err error
			headerID, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := common.HexToAddress("0x1234")
			toAddress := common.HexToAddress("0x5678")
			txHash := common.HexToHash("0x9876")
			txHashTwo := common.HexToHash("0x5432")
			txIndex := big.NewInt(123)
			transactions = []core.TransactionModel{{
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 0,
				GasPrice: 0,
				Hash:     txHash.Hex(),
				Nonce:    0,
				Raw:      []byte{},
				To:       toAddress.Hex(),
				TxIndex:  txIndex.Int64(),
				Value:    "0",
				Receipt:  core.Receipt{},
			}, {
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 1,
				GasPrice: 1,
				Hash:     txHashTwo.Hex(),
				Nonce:    1,
				Raw:      []byte{},
				To:       toAddress.Hex(),
				TxIndex:  1,
				Value:    "1",
			}}

			insertErr := repo.CreateTransactions(headerID, transactions)
			Expect(insertErr).NotTo(HaveOccurred())
		})

		It("adds transactions", func() {
			var dbTransactions []core.TransactionModel
			readErr := db.Select(&dbTransactions,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.transactions WHERE header_id = $1`, headerID)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(dbTransactions).To(ConsistOf(transactions))
		})

		It("silently ignores duplicate inserts", func() {
			insertTwoErr := repo.CreateTransactions(headerID, transactions)
			Expect(insertTwoErr).NotTo(HaveOccurred())

			var dbTransactions []core.TransactionModel
			readErr := db.Select(&dbTransactions,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.transactions WHERE header_id = $1`, headerID)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(len(dbTransactions)).To(Equal(2))
		})
	})

	Describe("creating a transaction in a sqlx tx", func() {
		It("adds a transaction", func() {
			headerID, err := repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := common.HexToAddress("0x1234")
			toAddress := common.HexToAddress("0x5678")
			txHash := common.HexToHash("0x9876")
			txIndex := big.NewInt(123)
			transaction := core.TransactionModel{
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 0,
				GasPrice: 0,
				Hash:     txHash.Hex(),
				Nonce:    0,
				Raw:      []byte{1, 2, 3},
				To:       toAddress.Hex(),
				TxIndex:  txIndex.Int64(),
				Value:    "0",
			}

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())
			_, insertErr := repo.CreateTransactionInTx(tx, headerID, transaction)
			Expect(insertErr).NotTo(HaveOccurred())
			commitErr := tx.Commit()
			Expect(commitErr).ToNot(HaveOccurred())

			var dbTransaction core.TransactionModel
			err = db.Get(&dbTransaction,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbTransaction).To(Equal(transaction))
		})

		It("silently upserts", func() {
			headerID, err := repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := common.HexToAddress("0x1234")
			toAddress := common.HexToAddress("0x5678")
			txHash := common.HexToHash("0x9876")
			txIndex := big.NewInt(123)
			transaction := core.TransactionModel{
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 0,
				GasPrice: 0,
				Hash:     txHash.Hex(),
				Nonce:    0,
				Raw:      []byte{},
				Receipt:  core.Receipt{},
				To:       toAddress.Hex(),
				TxIndex:  txIndex.Int64(),
				Value:    "0",
			}

			tx1Err := repo.CreateTransactions(headerID, []core.TransactionModel{transaction})
			Expect(tx1Err).NotTo(HaveOccurred())

			tx2Err := repo.CreateTransactions(headerID, []core.TransactionModel{transaction})
			Expect(tx2Err).NotTo(HaveOccurred())

			var dbTransactions []core.TransactionModel
			err = db.Select(&dbTransactions,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbTransactions)).To(Equal(1))
		})
	})

	Describe("Getting a header", func() {
		It("returns header if it exists", func() {
			_, createErr := repo.CreateOrUpdateHeader(header)
			Expect(createErr).NotTo(HaveOccurred())

			dbHeader, err := repo.GetHeader(header.BlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.Id).NotTo(BeZero())
			Expect(dbHeader.BlockNumber).To(Equal(header.BlockNumber))
			Expect(dbHeader.Hash).To(Equal(header.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(header.Raw))
			Expect(dbHeader.Timestamp).To(Equal(header.Timestamp))
		})

		It("returns header from any node", func() {
			_, createErr := repo.CreateOrUpdateHeader(header)
			Expect(createErr).NotTo(HaveOccurred())

			dbTwo := test_config.NewTestDB(test_config.NewTestNode())
			repoTwo := repositories.NewHeaderRepository(dbTwo)

			result, readErr := repoTwo.GetHeader(header.BlockNumber)

			Expect(readErr).NotTo(HaveOccurred())
			Expect(result.Raw).To(MatchJSON(header.Raw))
		})
	})

	Describe("Getting headers in range", func() {
		var blockTwo int64

		BeforeEach(func() {
			_, headerErrOne := repo.CreateOrUpdateHeader(header)
			Expect(headerErrOne).NotTo(HaveOccurred())
			blockTwo = header.BlockNumber + 1
			headerTwo := core.Header{
				BlockNumber: blockTwo,
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         header.Raw,
				Timestamp:   header.Timestamp,
			}
			_, headerErrTwo := repo.CreateOrUpdateHeader(headerTwo)
			Expect(headerErrTwo).NotTo(HaveOccurred())
		})

		It("returns all headers in block range", func() {
			dbHeaders, err := repo.GetHeadersInRange(header.BlockNumber, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(dbHeaders)).To(Equal(2))
		})

		It("does not return header outside of block range", func() {
			dbHeaders, err := repo.GetHeadersInRange(header.BlockNumber, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(dbHeaders)).To(Equal(1))
		})
	})

	Describe("Getting missing headers", func() {
		It("returns block numbers for headers not in the database", func() {
			_, createOneErr := repo.CreateOrUpdateHeader(fakes.GetFakeHeader(1))
			Expect(createOneErr).NotTo(HaveOccurred())

			_, createTwoErr := repo.CreateOrUpdateHeader(fakes.GetFakeHeader(3))
			Expect(createTwoErr).NotTo(HaveOccurred())

			_, createThreeErr := repo.CreateOrUpdateHeader(fakes.GetFakeHeader(5))
			Expect(createThreeErr).NotTo(HaveOccurred())

			missingBlockNumbers, err := repo.MissingBlockNumbers(1, 5)
			Expect(err).NotTo(HaveOccurred())

			Expect(missingBlockNumbers).To(ConsistOf([]int64{2, 4}))
		})

		It("treats headers created by _any_ node as not missing", func() {
			_, createOneErr := repo.CreateOrUpdateHeader(fakes.GetFakeHeader(1))
			Expect(createOneErr).NotTo(HaveOccurred())

			_, createTwoErr := repo.CreateOrUpdateHeader(fakes.GetFakeHeader(3))
			Expect(createTwoErr).NotTo(HaveOccurred())

			_, createThreeErr := repo.CreateOrUpdateHeader(fakes.GetFakeHeader(5))
			Expect(createThreeErr).NotTo(HaveOccurred())

			dbTwo := test_config.NewTestDB(test_config.NewTestNode())
			repoTwo := repositories.NewHeaderRepository(dbTwo)

			missingBlockNumbers, err := repoTwo.MissingBlockNumbers(1, 5)
			Expect(err).NotTo(HaveOccurred())

			Expect(missingBlockNumbers).To(ConsistOf([]int64{2, 4}))
		})
	})
})
