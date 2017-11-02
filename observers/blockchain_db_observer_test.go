package observers_test

import (
	"github.com/8thlight/vulcanizedb/config"
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/observers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Saving blocks to the database", func() {

	var db *sqlx.DB
	var err error

	BeforeEach(func() {
		pgConfig := config.DbConnectionString(config.Private().Database)
		db, err = sqlx.Connect("postgres", pgConfig)
		db.MustExec("DELETE FROM transactions")
		db.MustExec("DELETE FROM blocks")
	})

	AfterEach(func() {
		db.Close()
	})

	It("implements the observer interface", func() {
		var observer core.BlockchainObserver = observers.BlockchainDBObserver{Db: db}
		Expect(observer).NotTo(BeNil())
	})

	It("connects to the database", func() {
		Expect(err).Should(BeNil())
		Expect(db).ShouldNot(BeNil())
	})

	It("starts with no blocks", func() {
		var count int
		queryError := db.Get(&count, "SELECT COUNT(*) FROM blocks")
		Expect(queryError).Should(BeNil())
		Expect(count).Should(Equal(0))
	})

	It("inserts a block", func() {
		// setup a block in memory
		blockNumber := int64(123)
		gasLimit := int64(1000000)
		gasUsed := int64(10)
		blockHash := "x123"
		blockParentHash := "x456"
		blockNonce := "0x881db2ca900682e9a9"
		blockTime := int64(1508981640)
		uncleHash := "x789"
		blockSize := int64(1000)
		difficulty := int64(10)
		block := core.Block{
			Difficulty: difficulty,
			GasLimit:   gasLimit,
			GasUsed:    gasUsed,
			Hash:       blockHash,
			Nonce:      blockNonce,
			Number:     blockNumber,
			ParentHash: blockParentHash,
			Size:       blockSize,
			Time:       blockTime,
			UncleHash:  uncleHash,
		}

		// save the block to the database
		observer := observers.BlockchainDBObserver{Db: db}
		observer.NotifyBlockAdded(block)

		// find the saved block
		rows, err := db.Query(
			"SELECT  block_number, block_gaslimit, block_gasused, block_time, block_difficulty, block_hash, block_nonce, block_parenthash, block_size, uncle_hash FROM blocks")
		Expect(err).To(BeNil())
		var savedBlocks []core.Block
		for rows.Next() {
			var blockHash string
			var blockNonce string
			var blockNumber int64
			var blockParentHash string
			var blockSize int64
			var blockTime float64
			var difficulty int64
			var gasLimit float64
			var gasUsed float64
			var uncleHash string
			rows.Scan(&blockNumber, &gasLimit, &gasUsed, &blockTime, &difficulty, &blockHash, &blockNonce, &blockParentHash, &blockSize, &uncleHash)
			savedBlock := core.Block{
				Difficulty: difficulty,
				GasLimit:   int64(gasLimit),
				GasUsed:    int64(gasUsed),
				Hash:       blockHash,
				Nonce:      blockNonce,
				Number:     blockNumber,
				ParentHash: blockParentHash,
				Size:       blockSize,
				Time:       int64(blockTime),
				UncleHash:  uncleHash,
			}
			savedBlocks = append(savedBlocks, savedBlock)
		}
		// assert against the attributes
		Expect(len(savedBlocks)).To(Equal(1))
		Expect(savedBlocks[0].Difficulty).To(Equal(difficulty))
		Expect(savedBlocks[0].GasLimit).To(Equal(gasLimit))
		Expect(savedBlocks[0].GasUsed).To(Equal(gasUsed))
		Expect(savedBlocks[0].Hash).To(Equal(blockHash))
		Expect(savedBlocks[0].Nonce).To(Equal(blockNonce))
		Expect(savedBlocks[0].Number).To(Equal(blockNumber))
		Expect(savedBlocks[0].ParentHash).To(Equal(blockParentHash))
		Expect(savedBlocks[0].Size).To(Equal(blockSize))
		Expect(savedBlocks[0].Time).To(Equal(blockTime))
		Expect(savedBlocks[0].UncleHash).To(Equal(uncleHash))
	})

	var _ = Describe("Saving transactions to the database", func() {

		It("inserts a transaction", func() {
			gasLimit := int64(5000)
			gasPrice := int64(3)
			nonce := uint64(10000)
			to := "1234567890"
			value := int64(10)

			txRecord := core.Transaction{
				Hash:     "x1234",
				GasPrice: gasPrice,
				GasLimit: gasLimit,
				Nonce:    nonce,
				To:       to,
				Value:    value,
			}
			block := core.Block{Transactions: []core.Transaction{txRecord}}

			observer := observers.BlockchainDBObserver{Db: db}
			observer.NotifyBlockAdded(block)

			rows, err := db.Query("SELECT tx_hash, tx_nonce, tx_to, tx_gaslimit, tx_gasprice, tx_value FROM transactions")
			Expect(err).To(BeNil())

			var savedTransactions []core.Transaction
			for rows.Next() {
				var dbHash string
				var dbNonce uint64
				var dbTo string
				var dbGasLimit int64
				var dbGasPrice int64
				var dbValue int64
				rows.Scan(&dbHash, &dbNonce, &dbTo, &dbGasLimit, &dbGasPrice, &dbValue)
				savedTransaction := core.Transaction{
					GasLimit: dbGasLimit,
					GasPrice: dbGasPrice,
					Hash:     dbHash,
					Nonce:    dbNonce,
					To:       dbTo,
					Value:    dbValue,
				}
				savedTransactions = append(savedTransactions, savedTransaction)
			}

			Expect(len(savedTransactions)).To(Equal(1))
			savedTransaction := savedTransactions[0]
			Expect(savedTransaction.GasLimit).To(Equal(gasLimit))
			Expect(savedTransaction.GasPrice).To(Equal(gasPrice))
			Expect(savedTransaction.Hash).To(Equal(txRecord.Hash))
			Expect(savedTransaction.Nonce).To(Equal(nonce))
			Expect(savedTransaction.To).To(Equal(to))
			Expect(savedTransaction.Value).To(Equal(value))
		})

		It("associates the transaction with the block", func() {
			txRecord := core.Transaction{}
			block := core.Block{
				Transactions: []core.Transaction{txRecord},
			}

			observer := observers.BlockchainDBObserver{Db: db}
			observer.NotifyBlockAdded(block)

			blockRows, err := db.Query("SELECT id FROM blocks")
			Expect(err).To(BeNil())

			var actualBlockIds []int64
			for blockRows.Next() {
				var actualBlockId int64
				blockRows.Scan(&actualBlockId)
				actualBlockIds = append(actualBlockIds, actualBlockId)
			}

			transactionRows, err := db.Query("SELECT block_id FROM transactions")
			Expect(err).To(BeNil())

			var transactionBlockIds []int64
			for transactionRows.Next() {
				var transactionBlockId int64
				transactionRows.Scan(&transactionBlockId)
				transactionBlockIds = append(transactionBlockIds, transactionBlockId)
			}

			Expect(len(actualBlockIds)).To(Equal(1))
			Expect(len(transactionBlockIds)).To(Equal(1))
			Expect(transactionBlockIds[0]).To(Equal(actualBlockIds[0]))
		})
	})

})
