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
	"math/rand"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
)

type TransferLog struct {
	Id             int64  `db:"id"`
	VulvanizeLogId int64  `db:"vulcanize_log_id"`
	TokenName      string `db:"token_name"`
	Block          int64  `db:"block"`
	Tx             string `db:"tx"`
	From           string `db:"from_"`
	To             string `db:"to_"`
	Value          string `db:"value_"`
}

type NewOwnerLog struct {
	Id             int64  `db:"id"`
	VulvanizeLogId int64  `db:"vulcanize_log_id"`
	TokenName      string `db:"token_name"`
	Block          int64  `db:"block"`
	Tx             string `db:"tx"`
	Node           string `db:"node_"`
	Label          string `db:"label_"`
	Owner          string `db:"owner_"`
}

type LightTransferLog struct {
	Id        int64  `db:"id"`
	HeaderID  int64  `db:"header_id"`
	TokenName string `db:"token_name"`
	LogIndex  int64  `db:"log_idx"`
	TxIndex   int64  `db:"tx_idx"`
	From      string `db:"from_"`
	To        string `db:"to_"`
	Value     string `db:"value_"`
	RawLog    []byte `db:"raw_log"`
}

type LightNewOwnerLog struct {
	Id        int64  `db:"id"`
	HeaderID  int64  `db:"header_id"`
	TokenName string `db:"token_name"`
	LogIndex  int64  `db:"log_idx"`
	TxIndex   int64  `db:"tx_idx"`
	Node      string `db:"node_"`
	Label     string `db:"label_"`
	Owner     string `db:"owner_"`
	RawLog    []byte `db:"raw_log"`
}

type BalanceOf struct {
	Id        int64  `db:"id"`
	TokenName string `db:"token_name"`
	Block     int64  `db:"block"`
	Address   string `db:"who_"`
	Balance   string `db:"returned"`
}

type Resolver struct {
	Id        int64  `db:"id"`
	TokenName string `db:"token_name"`
	Block     int64  `db:"block"`
	Node      string `db:"node_"`
	Address   string `db:"returned"`
}

type Owner struct {
	Id        int64  `db:"id"`
	TokenName string `db:"token_name"`
	Block     int64  `db:"block"`
	Node      string `db:"node_"`
	Address   string `db:"returned"`
}

func SetupBC() core.BlockChain {
	infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
	rawRpcClient, err := rpc.Dial(infuraIPC)
	Expect(err).NotTo(HaveOccurred())
	rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
	ethClient := ethclient.NewClient(rawRpcClient)
	blockChainClient := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

	return blockChain
}

func SetupDBandBC() (*postgres.DB, core.BlockChain) {
	infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
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

func SetupTusdRepo(vulcanizeLogId *int64, wantedEvents, wantedMethods []string) (*postgres.DB, *contract.Contract) {
	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, core.Node{})
	Expect(err).NotTo(HaveOccurred())

	receiptRepository := repositories.ReceiptRepository{DB: db}
	logRepository := repositories.LogRepository{DB: db}
	blockRepository := *repositories.NewBlockRepository(db)

	blockNumber := rand.Int63()
	blockId := CreateBlock(blockNumber, blockRepository)

	receipts := []core.Receipt{{Logs: []core.Log{{}}}}

	err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)
	Expect(err).ToNot(HaveOccurred())

	err = logRepository.Get(vulcanizeLogId, `SELECT id FROM logs`)
	Expect(err).ToNot(HaveOccurred())

	info := SetupTusdContract(wantedEvents, wantedMethods)

	return db, info
}

func SetupTusdContract(wantedEvents, wantedMethods []string) *contract.Contract {
	p := mocks.NewParser(constants.TusdAbiString)
	err := p.Parse()
	Expect(err).ToNot(HaveOccurred())

	return contract.Contract{
		Name:          "TrueUSD",
		Address:       constants.TusdContractAddress,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		StartingBlock: 6194634,
		LastBlock:     6507323,
		Events:        p.GetEvents(wantedEvents),
		Methods:       p.GetSelectMethods(wantedMethods),
		MethodArgs:    map[string]bool{},
		FilterArgs:    map[string]bool{},
	}.Init()
}

func SetupENSRepo(vulcanizeLogId *int64, wantedEvents, wantedMethods []string) (*postgres.DB, *contract.Contract) {
	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, core.Node{})
	Expect(err).NotTo(HaveOccurred())

	receiptRepository := repositories.ReceiptRepository{DB: db}
	logRepository := repositories.LogRepository{DB: db}
	blockRepository := *repositories.NewBlockRepository(db)

	blockNumber := rand.Int63()
	blockId := CreateBlock(blockNumber, blockRepository)

	receipts := []core.Receipt{{Logs: []core.Log{{}}}}

	err = receiptRepository.CreateReceiptsAndLogs(blockId, receipts)
	Expect(err).ToNot(HaveOccurred())

	err = logRepository.Get(vulcanizeLogId, `SELECT id FROM logs`)
	Expect(err).ToNot(HaveOccurred())

	info := SetupENSContract(wantedEvents, wantedMethods)

	return db, info
}

func SetupENSContract(wantedEvents, wantedMethods []string) *contract.Contract {
	p := mocks.NewParser(constants.ENSAbiString)
	err := p.Parse()
	Expect(err).ToNot(HaveOccurred())

	return contract.Contract{
		Name:          "ENS-Registry",
		Address:       constants.EnsContractAddress,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		StartingBlock: 6194634,
		LastBlock:     6507323,
		Events:        p.GetEvents(wantedEvents),
		Methods:       p.GetSelectMethods(wantedMethods),
		MethodArgs:    map[string]bool{},
		FilterArgs:    map[string]bool{},
	}.Init()
}

func TearDown(db *postgres.DB) {
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM blocks`)
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

	_, err = tx.Exec(`DROP TABLE checked_headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE checked_headers (id SERIAL PRIMARY KEY, header_id INTEGER UNIQUE NOT NULL REFERENCES headers (id) ON DELETE CASCADE);`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP SCHEMA IF EXISTS full_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e CASCADE`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP SCHEMA IF EXISTS light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e CASCADE`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP SCHEMA IF EXISTS full_0x314159265dd8dbb310642f98f50c066173c1259b CASCADE`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP SCHEMA IF EXISTS light_0x314159265dd8dbb310642f98f50c066173c1259b CASCADE`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())
}

func CreateBlock(blockNumber int64, repository repositories.BlockRepository) (blockId int64) {
	blockId, err := repository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
	Expect(err).NotTo(HaveOccurred())

	return blockId
}
