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

	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Receipt Repository", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		receiptRepo repositories.ReceiptRepository
		headerRepo  datastore.HeaderRepository
		header      core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		receiptRepo = repositories.ReceiptRepository{}
		headerRepo = repositories.NewHeaderRepository(db)
		header = fakes.GetFakeHeader(rand.Int63())
	})

	Describe("creating a receipt", func() {
		It("adds a receipt in a tx", func() {
			headerID, err := headerRepo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := test_data.FakeAddress()
			toAddress := test_data.FakeAddress()
			txHash := test_data.FakeHash()
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

			receipt := core.Receipt{
				ContractAddress:   fromAddress.Hex(),
				TxHash:            txHash.Hex(),
				GasUsed:           uint64(rand.Int31()),
				CumulativeGasUsed: uint64(rand.Int31()),
				StateRoot:         test_data.FakeHash().Hex(),
				Rlp:               test_data.FakeHash().Bytes(),
			}

			_, receiptErr := receiptRepo.CreateReceiptInTx(headerID, txId, receipt, tx)
			Expect(receiptErr).ToNot(HaveOccurred())
			if receiptErr == nil {
				// this hangs if called when receiptsErr != nil
				commitErr := tx.Commit()
				Expect(commitErr).ToNot(HaveOccurred())
			} else {
				// lookup on public.receipts below hangs if tx is still open
				rollbackErr := tx.Rollback()
				Expect(rollbackErr).NotTo(HaveOccurred())
			}

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
			getAddressErr := db.Get(&addressId, `SELECT id FROM addresses WHERE address = $1`, fromAddress.Hex())
			Expect(getAddressErr).NotTo(HaveOccurred())

			var dbReceipt idModel
			getReceiptErr := db.Get(&dbReceipt,
				`SELECT transaction_id, contract_address_id, cumulative_gas_used, gas_used, state_root, status, tx_hash, rlp
				FROM public.receipts WHERE header_id = $1`, headerID)
			Expect(getReceiptErr).NotTo(HaveOccurred())

			Expect(dbReceipt.TransactionId).To(Equal(txId))
			Expect(dbReceipt.TxHash).To(Equal(txHash.Hex()))
			Expect(dbReceipt.ContractAddressId).To(Equal(addressId))
			Expect(dbReceipt.CumulativeGasUsed).To(Equal(receipt.CumulativeGasUsed))
			Expect(dbReceipt.GasUsed).To(Equal(receipt.GasUsed))
			Expect(dbReceipt.StateRoot).To(Equal(receipt.StateRoot))
			Expect(dbReceipt.Status).To(Equal(0))
			Expect(dbReceipt.Rlp).To(Equal(receipt.Rlp))
		})
	})
})
