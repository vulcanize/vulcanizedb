package repository_test

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

var _ = Describe("receipt updated trigger", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		receiptRepo repositories.ReceiptRepository
		headerRepo  datastore.HeaderRepository
		headerID    int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		receiptRepo = repositories.ReceiptRepository{}
		headerRepo = repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepo.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	type dbEvent struct {
		Created string
		Updated string
	}

	It("indicates when a receipt record was created or updated", func() {
		var receiptRes dbEvent
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

		createdReceiptErr := db.Get(&receiptRes, `SELECT created, updated FROM receipts`)
		Expect(createdReceiptErr).NotTo(HaveOccurred())
		Expect(receiptRes.Created).To(Equal(receiptRes.Updated))

		_, updateErr := db.Exec(`UPDATE public.receipts SET tx_hash = '{"new_hash"}' WHERE header_id = $1`, headerID)
		Expect(updateErr).NotTo(HaveOccurred())
		updatedReceiptErr := db.Get(&receiptRes, `SELECT created, updated FROM receipts`)
		Expect(updatedReceiptErr).NotTo(HaveOccurred())
		Expect(receiptRes.Created).NotTo(Equal(receiptRes.Updated))
	})
})
