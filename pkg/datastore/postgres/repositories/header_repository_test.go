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
	"database/sql"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Block header repository", func() {
	var (
		rawHeader []byte
		err       error
		timestamp string
		db        *postgres.DB
		repo      repositories.HeaderRepository
		header    core.Header
	)

	BeforeEach(func() {
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
		timestamp = big.NewInt(123456789).String()

		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = repositories.NewHeaderRepository(db)
		header = core.Header{
			BlockNumber: 100,
			Bloom:       "0x0000",
			Hash:        common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex(),
			Raw:         rawHeader,
			Timestamp:   timestamp,
		}
	})

	Describe("creating or updating a header", func() {
		It("adds a header", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			var dbHeader core.Header
			err = db.Get(&dbHeader, `SELECT block_number, bloom, hash, raw, block_timestamp FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.BlockNumber).To(Equal(header.BlockNumber))
			Expect(dbHeader.Bloom).To(Equal(header.Bloom))
			Expect(dbHeader.Hash).To(Equal(header.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(header.Raw))
			Expect(dbHeader.Timestamp).To(Equal(header.Timestamp))
		})

		It("adds node data to header", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			var ethNodeId int64
			err = db.Get(&ethNodeId, `SELECT eth_node_id FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(ethNodeId).To(Equal(db.NodeID))
			var ethNodeFingerprint string
			err = db.Get(&ethNodeFingerprint, `SELECT eth_node_fingerprint FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(ethNodeFingerprint).To(Equal(db.Node.ID))
		})

		It("returns valid header exists error if attempting duplicate headers", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(repositories.ErrValidHeaderExists))

			var dbHeaders []core.Header
			err = db.Select(&dbHeaders, `SELECT block_number, hash, raw FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbHeaders)).To(Equal(1))
		})

		It("replaces header if hash is different", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			headerTwo := core.Header{
				BlockNumber: header.BlockNumber,
				Bloom:       "0x0001",
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}

			_, err = repo.CreateOrUpdateHeader(headerTwo)

			Expect(err).NotTo(HaveOccurred())
			var dbHeader core.Header
			err = db.Get(&dbHeader, `SELECT block_number, bloom, hash, raw FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.Hash).To(Equal(headerTwo.Hash))
			Expect(dbHeader.Bloom).To(Equal(headerTwo.Bloom))
			Expect(dbHeader.Raw).To(MatchJSON(headerTwo.Raw))
		})

		It("does not replace header if node fingerprint is different", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())

			repoTwo := repositories.NewHeaderRepository(dbTwo)
			headerTwo := core.Header{
				BlockNumber: header.BlockNumber,
				Bloom:       "0x0000",
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}

			_, err = repoTwo.CreateOrUpdateHeader(headerTwo)

			Expect(err).NotTo(HaveOccurred())
			var dbHeaders []core.Header
			err = dbTwo.Select(&dbHeaders, `SELECT block_number, bloom, hash, raw FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbHeaders)).To(Equal(2))
		})

		It("only replaces header with matching node fingerprint", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())

			repoTwo := repositories.NewHeaderRepository(dbTwo)
			headerTwo := core.Header{
				BlockNumber: header.BlockNumber,
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}
			_, err = repoTwo.CreateOrUpdateHeader(headerTwo)
			headerThree := core.Header{
				BlockNumber: header.BlockNumber,
				Hash:        common.BytesToHash([]byte{1, 1, 1, 1, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}

			_, err = repoTwo.CreateOrUpdateHeader(headerThree)

			Expect(err).NotTo(HaveOccurred())
			var dbHeaders []core.Header
			err = dbTwo.Select(&dbHeaders, `SELECT block_number, hash, raw FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbHeaders)).To(Equal(2))
			Expect(dbHeaders[0].Hash).To(Or(Equal(header.Hash), Equal(headerThree.Hash)))
			Expect(dbHeaders[1].Hash).To(Or(Equal(header.Hash), Equal(headerThree.Hash)))
			Expect(dbHeaders[0].Raw).To(Or(MatchJSON(header.Raw), MatchJSON(headerThree.Raw)))
			Expect(dbHeaders[1].Raw).To(Or(MatchJSON(header.Raw), MatchJSON(headerThree.Raw)))
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

			receiptRepo := repositories.HeaderSyncReceiptRepository{}
			_, receiptErr := receiptRepo.CreateHeaderSyncReceiptInTx(headerID, txId, receipt, tx)
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
				FROM public.header_sync_receipts WHERE header_id = $1`, headerID)
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
			err = db.Select(&dbTransactions,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.header_sync_transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbTransactions).To(ConsistOf(transactions))
		})

		It("silently ignores duplicate inserts", func() {
			insertTwoErr := repo.CreateTransactions(headerID, transactions)
			Expect(insertTwoErr).NotTo(HaveOccurred())

			var dbTransactions []core.TransactionModel
			err = db.Select(&dbTransactions,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.header_sync_transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
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
			commitErr := tx.Commit()
			Expect(commitErr).ToNot(HaveOccurred())
			Expect(insertErr).NotTo(HaveOccurred())

			var dbTransaction core.TransactionModel
			err = db.Get(&dbTransaction,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.header_sync_transactions WHERE header_id = $1`, headerID)
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

			tx1, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())
			txId1, insertErr := repo.CreateTransactionInTx(tx1, headerID, transaction)
			commit1Err := tx1.Commit()
			Expect(commit1Err).ToNot(HaveOccurred())
			Expect(insertErr).NotTo(HaveOccurred())

			tx2, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())
			txId2, insertErr := repo.CreateTransactionInTx(tx2, headerID, transaction)
			commit2Err := tx2.Commit()
			Expect(commit2Err).ToNot(HaveOccurred())
			Expect(insertErr).NotTo(HaveOccurred())
			Expect(txId1).To(Equal(txId2))

			var dbTransactions []core.TransactionModel
			err = db.Select(&dbTransactions,
				`SELECT hash, gas_limit, gas_price, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.header_sync_transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbTransactions)).To(Equal(1))
		})
	})

	Describe("Getting a header", func() {
		It("returns header if it exists", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			dbHeader, err := repo.GetHeader(header.BlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.Id).NotTo(BeZero())
			Expect(dbHeader.BlockNumber).To(Equal(header.BlockNumber))
			Expect(dbHeader.Bloom).To(Equal(header.Bloom))
			Expect(dbHeader.Hash).To(Equal(header.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(header.Raw))
			Expect(dbHeader.Timestamp).To(Equal(header.Timestamp))
		})

		It("does not return header for a different node fingerprint", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())
			repoTwo := repositories.NewHeaderRepository(dbTwo)

			_, err = repoTwo.GetHeader(header.BlockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("Getting missing headers", func() {
		It("returns block numbers for headers not in the database", func() {
			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 1,
				Bloom:       "0x0000",
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 3,
				Bloom:       "0x0000",
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 5,
				Bloom:       "0x0000",
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			missingBlockNumbers, err := repo.MissingBlockNumbers(1, 5, db.Node.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(missingBlockNumbers).To(ConsistOf([]int64{2, 4}))
		})

		It("does not count headers created by a different node fingerprint", func() {
			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 1,
				Bloom:       "0x0000",
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 3,
				Bloom:       "0x0000",
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 5,
				Bloom:       "0x0000",
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())
			repoTwo := repositories.NewHeaderRepository(dbTwo)

			missingBlockNumbers, err := repoTwo.MissingBlockNumbers(1, 5, nodeTwo.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(missingBlockNumbers).To(ConsistOf([]int64{1, 2, 3, 4, 5}))
		})
	})
})
