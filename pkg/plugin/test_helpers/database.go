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

package test_helpers

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

func SetupDBandBC() (*postgres.DB, core.BlockChain) {
	infuraIPC := "http://kovan0.vulcanize.io:8545"
	rawRpcClient, err := rpc.Dial(infuraIPC)
	Expect(err).NotTo(HaveOccurred())
	rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
	ethClient := ethclient.NewClient(rawRpcClient)
	blockChainClient := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, blockChain.Node())
	Expect(err).NotTo(HaveOccurred())

	return db, blockChain
}

func TearDown(db *postgres.DB) {
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM log_filters`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM transactions`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM receipts`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM checked_headers`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())
}

func DropTestSchema(db *postgres.DB, schema string) {
	_, err := db.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, schema))
	Expect(err).NotTo(HaveOccurred())
}
