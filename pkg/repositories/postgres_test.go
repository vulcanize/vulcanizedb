package repositories_test

import (
	"fmt"
	"strings"

	"io/ioutil"
	"log"

	"math/big"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var _ = Describe("Postgres repository", func() {
	var repository repositories.Postgres

	It("connects to the database", func() {
		cfg, _ := config.NewConfig("private")
		pgConfig := config.DbConnectionString(cfg.Database)
		db, err := sqlx.Connect("postgres", pgConfig)
		Expect(err).Should(BeNil())
		Expect(db).ShouldNot(BeNil())
	})

	BeforeEach(func() {
		node := core.Node{
			GenesisBlock: "GENESIS",
			NetworkId:    1,
			Id:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		cfg, _ := config.NewConfig("private")
		repository, _ = repositories.NewPostgres(cfg.Database, node)
		repositories.ClearData(repository)
	})

	It("serializes big.Int to db", func() {
		// postgres driver doesn't support go big.Int type
		// various casts in golang uint64, int64, overflow for
		// transaction value (in wei) even though
		// postgres numeric can handle an arbitrary
		// sized int, so use string representation of big.Int
		// and cast on insert

		cfg, _ := config.NewConfig("private")
		pgConfig := config.DbConnectionString(cfg.Database)
		db, err := sqlx.Connect("postgres", pgConfig)

		bi := new(big.Int)
		bi.SetString("34940183920000000000", 10)
		Expect(bi.String()).To(Equal("34940183920000000000"))

		defer db.Exec(`DROP TABLE IF EXISTS example`)
		_, err = db.Exec("CREATE TABLE example ( id INTEGER, data NUMERIC )")
		Expect(err).ToNot(HaveOccurred())

		sqlStatement := `  
			INSERT INTO example (id, data)
			VALUES (1, cast($1 AS NUMERIC))`
		_, err = db.Exec(sqlStatement, bi.String())
		Expect(err).ToNot(HaveOccurred())

		var data string
		err = db.QueryRow(`SELECT data FROM example WHERE id = 1`).Scan(&data)
		Expect(err).ToNot(HaveOccurred())

		Expect(bi.String()).To(Equal(data))
		actual := new(big.Int)
		actual.SetString(data, 10)
		Expect(actual).To(Equal(bi))
	})

	It("does not commit block if block is invalid", func() {
		//badNonce violates db Nonce field length
		badNonce := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		badBlock := core.Block{
			Number:       123,
			Nonce:        badNonce,
			Transactions: []core.Transaction{},
		}
		cfg, _ := config.NewConfig("private")
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1, Id: "x123", ClientName: "geth"}
		repository, _ := repositories.NewPostgres(cfg.Database, node)

		err1 := repository.CreateOrUpdateBlock(badBlock)
		savedBlock, err2 := repository.FindBlockByNumber(123)

		Expect(err1).To(HaveOccurred())
		Expect(err2).To(HaveOccurred())
		Expect(savedBlock).To(BeZero())
	})

	It("throws error when can't connect to the database", func() {
		invalidDatabase := config.Database{}
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1, Id: "x123", ClientName: "geth"}
		_, err := repositories.NewPostgres(invalidDatabase, node)
		Expect(err).To(Equal(repositories.ErrDBConnectionFailed))
	})

	It("throws error when can't create node", func() {
		cfg, _ := config.NewConfig("private")
		badHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		node := core.Node{GenesisBlock: badHash, NetworkId: 1, Id: "x123", ClientName: "geth"}
		_, err := repositories.NewPostgres(cfg.Database, node)
		Expect(err).To(Equal(repositories.ErrUnableToSetNode))
	})

	It("does not commit log if log is invalid", func() {
		//badTxHash violates db tx_hash field length
		badTxHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		badLog := core.Log{
			Address:     "x123",
			BlockNumber: 1,
			TxHash:      badTxHash,
		}
		cfg, _ := config.NewConfig("private")
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1, Id: "x123", ClientName: "geth"}
		repository, _ := repositories.NewPostgres(cfg.Database, node)

		err := repository.CreateLogs([]core.Log{badLog})
		savedBlock := repository.FindLogs("x123", 1)

		Expect(err).ToNot(BeNil())
		Expect(savedBlock).To(BeNil())
	})

	It("does not commit block or transactions if transaction is invalid", func() {
		//badHash violates db To field length
		badHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		badTransaction := core.Transaction{To: badHash}
		block := core.Block{
			Number:       123,
			Transactions: []core.Transaction{badTransaction},
		}
		cfg, _ := config.NewConfig("private")
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1, Id: "x123", ClientName: "geth"}
		repository, _ := repositories.NewPostgres(cfg.Database, node)

		err1 := repository.CreateOrUpdateBlock(block)
		savedBlock, err2 := repository.FindBlockByNumber(123)

		Expect(err1).To(HaveOccurred())
		Expect(err2).To(HaveOccurred())
		Expect(savedBlock).To(BeZero())
	})

})
