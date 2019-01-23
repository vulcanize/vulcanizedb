// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package postgres_test

import (
	"fmt"
	"strings"

	"math/big"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Postgres DB", func() {
	var sqlxdb *sqlx.DB

	It("connects to the database", func() {
		var err error
		pgConfig := config.DbConnectionString(test_config.DBConfig)

		sqlxdb, err = sqlx.Connect("postgres", pgConfig)

		Expect(err).Should(BeNil())
		Expect(sqlxdb).ShouldNot(BeNil())
	})

	It("serializes big.Int to db", func() {
		// postgres driver doesn't support go big.Int type
		// various casts in golang uint64, int64, overflow for
		// transaction value (in wei) even though
		// postgres numeric can handle an arbitrary
		// sized int, so use string representation of big.Int
		// and cast on insert

		pgConnectString := config.DbConnectionString(test_config.DBConfig)
		db, err := sqlx.Connect("postgres", pgConnectString)
		Expect(err).NotTo(HaveOccurred())

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
		node := core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "geth"}
		db := test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		blocksRepository := repositories.NewBlockRepository(db)

		_, err1 := blocksRepository.CreateOrUpdateBlock(badBlock)

		Expect(err1).To(HaveOccurred())
		savedBlock, err2 := blocksRepository.GetBlock(123)
		Expect(err2).To(HaveOccurred())
		Expect(savedBlock).To(BeZero())
	})

	It("throws error when can't connect to the database", func() {
		invalidDatabase := config.Database{}
		node := core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "geth"}

		_, err := postgres.NewDB(invalidDatabase, node)

		Expect(err).To(Equal(postgres.ErrDBConnectionFailed))
	})

	It("throws error when can't create node", func() {
		badHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		node := core.Node{GenesisBlock: badHash, NetworkID: 1, ID: "x123", ClientName: "geth"}
		_, err := postgres.NewDB(test_config.DBConfig, node)
		Expect(err).To(Equal(postgres.ErrUnableToSetNode))
	})

	It("does not commit log if log is invalid", func() {
		//badTxHash violates db tx_hash field length
		badTxHash := fmt.Sprintf("x %s", strings.Repeat("1", 100))
		badLog := core.Log{
			Address:     "x123",
			BlockNumber: 1,
			TxHash:      badTxHash,
		}
		node := core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "geth"}
		db, _ := postgres.NewDB(test_config.DBConfig, node)
		logRepository := repositories.LogRepository{DB: db}

		err := logRepository.CreateLogs([]core.Log{badLog}, 123)

		Expect(err).ToNot(BeNil())
		savedBlock := logRepository.GetLogs("x123", 1)
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
		node := core.Node{GenesisBlock: "GENESIS", NetworkID: 1, ID: "x123", ClientName: "geth"}
		db, _ := postgres.NewDB(test_config.DBConfig, node)
		blockRepository := repositories.NewBlockRepository(db)

		_, err1 := blockRepository.CreateOrUpdateBlock(block)

		Expect(err1).To(HaveOccurred())
		savedBlock, err2 := blockRepository.GetBlock(123)
		Expect(err2).To(HaveOccurred())
		Expect(savedBlock).To(BeZero())
	})
})
