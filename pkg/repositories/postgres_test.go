package repositories_test

import (
	"fmt"
	"strings"

	"io/ioutil"
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/testing"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var _ = Describe("Postgres repository", func() {

	It("connects to the database", func() {
		cfg, _ := config.NewConfig("private")
		pgConfig := config.DbConnectionString(cfg.Database)
		db, err := sqlx.Connect("postgres", pgConfig)
		Expect(err).Should(BeNil())
		Expect(db).ShouldNot(BeNil())
	})

	testing.AssertRepositoryBehavior(func(node core.Node) repositories.Repository {
		cfg, _ := config.NewConfig("private")
		repository, _ := repositories.NewPostgres(cfg.Database, node)
		testing.ClearData(repository)
		return repository
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
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1}
		repository, _ := repositories.NewPostgres(cfg.Database, node)

		err1 := repository.CreateOrUpdateBlock(badBlock)
		savedBlock, err2 := repository.FindBlockByNumber(123)

		Expect(err1).To(HaveOccurred())
		Expect(err2).To(HaveOccurred())
		Expect(savedBlock).To(BeZero())
	})

	It("throws error when can't connect to the database", func() {
		invalidDatabase := config.Database{}
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1}
		_, err := repositories.NewPostgres(invalidDatabase, node)
		Expect(err).To(Equal(repositories.ErrDBConnectionFailed))
	})

	It("throws error when can't create node", func() {
		cfg, _ := config.NewConfig("private")
		badHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		node := core.Node{GenesisBlock: badHash, NetworkId: 1}
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
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1}
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
		node := core.Node{GenesisBlock: "GENESIS", NetworkId: 1}
		repository, _ := repositories.NewPostgres(cfg.Database, node)

		err1 := repository.CreateOrUpdateBlock(block)
		savedBlock, err2 := repository.FindBlockByNumber(123)

		Expect(err1).To(HaveOccurred())
		Expect(err2).To(HaveOccurred())
		Expect(savedBlock).To(BeZero())
	})

})
