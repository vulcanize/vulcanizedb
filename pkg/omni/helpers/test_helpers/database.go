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
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/parser"
)

var ExpectedTransferFilter = filters.LogFilter{
	Name:      "Transfer",
	Address:   constants.TusdContractAddress,
	ToBlock:   -1,
	FromBlock: 5197514,
	Topics:    core.Topics{constants.TransferEvent.Signature()},
}

var ExpectedApprovalFilter = filters.LogFilter{
	Name:      "Approval",
	Address:   constants.TusdContractAddress,
	ToBlock:   -1,
	FromBlock: 5197514,
	Topics:    core.Topics{constants.ApprovalEvent.Signature()},
}

type TransferLog struct {
	Id             int64  `db:"id"`
	VulvanizeLogId int64  `db:"vulcanize_log_id"`
	TokenName      string `db:"token_name"`
	TokenAddress   string `db:"token_address"`
	EventName      string `db:"event_name"`
	Block          int64  `db:"block"`
	Tx             string `db:"tx"`
	From           string `db:"from_"`
	To             string `db:"to_"`
	Value          string `db:"value_"`
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
	blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)

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
	blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)

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
	p := parser.NewParser("")
	err := p.Parse(constants.TusdContractAddress)
	Expect(err).ToNot(HaveOccurred())

	return &contract.Contract{
		Name:           "TrueUSD",
		Address:        constants.TusdContractAddress,
		Abi:            p.Abi(),
		ParsedAbi:      p.ParsedAbi(),
		StartingBlock:  5197514,
		LastBlock:      6507323,
		Events:         p.GetEvents(wantedEvents),
		Methods:        p.GetMethods(wantedMethods),
		EventAddrs:     map[string]bool{},
		MethodAddrs:    map[string]bool{},
		TknHolderAddrs: map[string]bool{},
	}
}

func TearDown(db *postgres.DB) {
	_, err := db.Query(`DELETE FROM blocks`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Query(`DELETE FROM logs`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Query(`DELETE FROM transactions`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Query(`DELETE FROM receipts`)
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Query(`DROP SCHEMA IF EXISTS c0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E CASCADE`)
	Expect(err).NotTo(HaveOccurred())
}

func CreateBlock(blockNumber int64, repository repositories.BlockRepository) (blockId int64) {
	blockId, err := repository.CreateOrUpdateBlock(core.Block{Number: blockNumber})
	Expect(err).NotTo(HaveOccurred())

	return blockId
}
