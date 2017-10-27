package core_test

import (
	"math/big"

	"fmt"

	"github.com/8thlight/vulcanizedb/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "vulcanize"
)

var _ = Describe("Saving blocks to the database", func() {

	var db *sqlx.DB
	var err error
	var gethTransaction *types.Transaction

	BeforeEach(func() {

		blockName := []byte("0x28f9a8d33109c87bda4a9ea890792421c710fe1c")
		addr := common.BytesToAddress(blockName)
		nonce := uint64(18848)
		amt := big.NewInt(0)
		gasLimit := big.NewInt(0)
		gasPrice := big.NewInt(0)
		data := []byte{}
		gethTransaction = types.NewTransaction(nonce, addr, amt, gasLimit, gasPrice, data)

		pgConfig := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err = sqlx.Connect("postgres", pgConfig)
		db.MustExec("DELETE FROM blocks")
	})

	It("implements the observer interface", func() {
		var observer core.BlockchainObserver = core.BlockchainDBObserver{Db: db}
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
		blockNumber := big.NewInt(1)
		gasLimit := big.NewInt(1000000)
		gasUsed := big.NewInt(10)
		blockTime := big.NewInt(1508981640)
		block := core.Block{Number: blockNumber, GasLimit: gasLimit, GasUsed: gasUsed, Time: blockTime}

		// save the block to the database
		observer := core.BlockchainDBObserver{Db: db}
		observer.NotifyBlockAdded(block)

		// find the saved block
		rows, err := db.Query("SELECT block_number, block_gaslimit, block_gasused, block_time FROM blocks")
		Expect(err).To(BeNil())
		var savedBlocks []core.Block
		for rows.Next() {
			var blockNumber int64
			var blockTime float64
			var gasLimit float64
			var gasUsed float64
			rows.Scan(&blockNumber, &gasLimit, &gasUsed, &blockTime)
			savedBlock := core.Block{
				GasUsed: big.NewInt(int64(gasUsed)),
				GasLimit: big.NewInt(int64(gasLimit)),
				Number:   big.NewInt(blockNumber),
				Time:   big.NewInt(int64(blockTime)),
			}
			savedBlocks = append(savedBlocks, savedBlock)
		}
		// assert against the attributes
		Expect(len(savedBlocks)).To(Equal(1))
		Expect(savedBlocks[0].Number.Int64()).To(Equal(blockNumber.Int64()))
		Expect(savedBlocks[0].GasLimit.Int64()).To(Equal(gasLimit.Int64()))
		Expect(savedBlocks[0].GasUsed.Int64()).To(Equal(gasUsed.Int64()))
		Expect(savedBlocks[0].Time).To(Equal(blockTime))
	})

})
