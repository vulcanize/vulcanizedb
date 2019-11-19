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
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/test_config"
)

var _ = Describe("Header Sync Receipt Repo", func() {
	var (
		rawHeader   []byte
		err         error
		timestamp   string
		db          *postgres.DB
		receiptRepo repositories.HeaderSyncReceiptRepository
		headerRepo  repositories.HeaderRepository
		header      core.Header
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
