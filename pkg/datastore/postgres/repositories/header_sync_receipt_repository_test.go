// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repositories_test

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/big"
)

var _ = Describe("Header Sync Receipt Repo", func() {
	var (
		rawHeader []byte
		err       error
		timestamp string
		db        *postgres.DB
		receiptRepo repositories.HeaderSyncReceiptRepository
		headerRepo repositories.HeaderRepository
		header    core.Header
	)

	BeforeEach(func() {
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
		timestamp = big.NewInt(123456789).String()

		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		receiptRepo = repositories.HeaderSyncReceiptRepository{}
		headerRepo = repositories.NewHeaderRepository(db)
		header = core.Header{
			BlockNumber: 100,
			Hash:        common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex(),
			Raw:         rawHeader,
			Timestamp:   timestamp,
		}
	})
	Describe("creating a receipt", func() {
		It("adds a receipt in a tx", func() {
			headerID, err := headerRepo.CreateOrUpdateHeader(header)
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
			txId, txErr := headerRepo.CreateTransactionInTx(tx, headerID, transaction)
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
})
