package observers_test

import (
	"runtime"

	"github.com/8thlight/vulcanizedb/config"
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/observers"
	"github.com/8thlight/vulcanizedb/repositories"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	_, filename, _, _ = runtime.Caller(0)
)

var _ = Describe("Saving blocks to the database", func() {

	var db *sqlx.DB
	var err error

	BeforeEach(func() {
		cfg := config.NewConfig("private")
		pgConfig := config.DbConnectionString(cfg.Database)
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

	It("saves a block with one transaction", func() {
		block := core.Block{
			Number:       123,
			Transactions: []core.Transaction{{}},
		}

		observer := observers.BlockchainDBObserver{Db: db}
		observer.NotifyBlockAdded(block)

		repository := repositories.NewPostgres(db)
		savedBlock := repository.FindBlockByNumber(123)
		Expect(savedBlock).NotTo(BeNil())
		Expect(len(savedBlock.Transactions)).To(Equal(1))
	})

})
