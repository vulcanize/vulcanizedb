package repositories_test

import (
	"fmt"
	"strings"

	"github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/8thlight/vulcanizedb/pkg/repositories/testing"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Postgres repository", func() {

	It("connects to the database", func() {
		cfg, _ := config.NewConfig("private")
		pgConfig := config.DbConnectionString(cfg.Database)
		db, err := sqlx.Connect("postgres", pgConfig)
		Expect(err).Should(BeNil())
		Expect(db).ShouldNot(BeNil())
	})

	testing.AssertRepositoryBehavior(func() repositories.Repository {
		cfg, _ := config.NewConfig("private")
		repository := repositories.NewPostgres(cfg.Database)
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
		repository := repositories.NewPostgres(cfg.Database)

		err := repository.CreateBlock(badBlock)
		savedBlock := repository.FindBlockByNumber(123)

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
		repository := repositories.NewPostgres(cfg.Database)

		err := repository.CreateBlock(block)
		savedBlock := repository.FindBlockByNumber(123)

		Expect(err).ToNot(BeNil())
		Expect(savedBlock).To(BeNil())
	})

})
