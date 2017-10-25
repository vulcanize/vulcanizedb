package core_test

import (
	"math/big"

	"fmt"

	"github.com/8thlight/vulcanizedb/core"
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

	BeforeEach(func() {
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
		block := core.Block{Number: big.NewInt(1)}

		observer := core.BlockchainDBObserver{Db: db}
		observer.NotifyBlockAdded(block)

		rows, err := db.Queryx("SELECT * FROM blocks")
		Expect(err).To(BeNil())
		var savedBlocks []core.BlockRecord
		for rows.Next() {
			var savedBlock core.BlockRecord
			rows.StructScan(&savedBlock)
			savedBlocks = append(savedBlocks, savedBlock)
		}

		Expect(len(savedBlocks)).To(Equal(1))
		Expect(savedBlocks[0].BlockNumber).To(Equal(int64(1)))
	})

})
