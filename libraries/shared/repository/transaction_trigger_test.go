package repository_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/big"
)

var _ = Describe("transaction updated trigger", func() {
	var (
		db       = test_config.NewTestDB(test_config.NewTestNode())
		headerID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	type dbTrans struct {
		Created string
		Updated string
	}

	It("indicates when a record was created or updated", func() {
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		fromAddress := common.HexToAddress("0x1234")
		toAddress := common.HexToAddress("0x5678")
		txHash := common.HexToHash("0x9876")
		txIndex := big.NewInt(123)
		transaction := []core.TransactionModel{{
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
		}}

		txError := headerRepository.CreateTransactions(headerID, transaction)
		Expect(txError).NotTo(HaveOccurred())

		var transUpdateRes dbTrans
		initialTxErr := db.Get(&transUpdateRes, `SELECT created, updated FROM public.transactions`)
		Expect(initialTxErr).NotTo(HaveOccurred())
		Expect(transUpdateRes.Created).To(Equal(transUpdateRes.Updated))

		_, updateErr := db.Exec(`UPDATE public.transactions SET hash = '{"new_hash"}' WHERE header_id = $1`, headerID)
		Expect(updateErr).NotTo(HaveOccurred())
		updatedTransErr := db.Get(&transUpdateRes, `SELECT created, updated FROM public.transactions`)
		Expect(updatedTransErr).NotTo(HaveOccurred())
		Expect(transUpdateRes.Created).NotTo(Equal(transUpdateRes.Updated))
	})
})
